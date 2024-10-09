package service

import (
	"crypto/rand"
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/session"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

type ClusterManager struct {
	RemoteService
	ContainerService
}

func (c *ClusterManager) QueryCluster(ctx *gin.Context, sessions *session.Manager) {
	//获取查询参数
	db := database.DB
	clusters := make([]model.Cluster, 0)
	session := sessions.Start(ctx.Writer, ctx.Request)
	role := session.Get("role")
	userId := session.Get("userId")
	log.Println("role", role, "userId", userId)
	if role == "admin" {
		if err := db.Find(&clusters).Error; err != nil {
			fmt.Println(err.Error())
		}
	} else {
		db2 := database.DB
		userClusterList := []model.UserCluster{}
		db2.Where("user_id = ?", userId).Find(&userClusterList)
		ids := []int{}
		for i := range userClusterList {
			ids = append(ids, userClusterList[i].ClusterId)
		}
		db.Where("id in ?", ids).Find(&clusters)
	}

	var clusterVOList = make([]vo.ClusterVO, len(clusters))
	var wg sync.WaitGroup
	wg.Add(len(clusters))
	for i, cluster := range clusters {
		go func(cluster model.Cluster, i int) {
			clusterVO := vo.ClusterVO{
				Name:            cluster.Name,
				ClusterName:     cluster.ClusterName,
				Description:     cluster.Description,
				ID:              cluster.ID,
				CreatedAt:       cluster.CreatedAt,
				UpdatedAt:       cluster.UpdatedAt,
				Ip:              cluster.Ip,
				Port:            cluster.Port,
				Username:        cluster.Username,
				ClusterPassword: cluster.Password,
			}
			clusterVO.GameArchive = c.GetRemoteGameArchive(cluster)
			clusterVO.Status = c.GetRemoteLevelStatus(cluster)
			clusterVOList[i] = clusterVO
			wg.Done()
		}(cluster, i)
	}
	wg.Wait()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterVOList,
	})

}

func (c *ClusterManager) CreateCluster(cluster *model.Cluster) {

	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	containerId, err := c.CreateContainer(*cluster)
	cluster.ContainerId = containerId
	cluster.Uuid = containerId
	err = db.Create(&cluster).Error

	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry" {
			log.Panicln("唯一索引冲突！", err)
		}
		log.Panicln("创建房间失败！", err)
	}
	tx.Commit()

}

func (c *ClusterManager) UpdateCluster(cluster *model.Cluster) {
	db := database.DB
	oldCluster := &model.Cluster{}
	db.Where("ID = ?", cluster.ID).First(oldCluster)
	if oldCluster.ClusterName == "" {
		log.Panicln("未找到当前存档 clusterName: ", cluster.ClusterName, cluster.ID)
	}
	oldCluster.Name = cluster.Name
	db.Updates(oldCluster)

}

func (c *ClusterManager) DeleteCluster(clusterName string) (*model.Cluster, error) {

	if clusterName == "" {
		log.Panicln("cluster is not allow null")
	}

	db := database.DB
	cluster := model.Cluster{}
	result := db.Where("cluster_name = ?", clusterName).Unscoped().Delete(&cluster)

	err := c.DeleteContainer(cluster.ContainerId)
	if err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &cluster, nil
}

func (c *ClusterManager) FindClusterByUuid(uuid string) *model.Cluster {
	db := database.DB
	cluster := &model.Cluster{}
	db.Where("uuid=?", uuid).First(cluster)
	return cluster
}

// 生成随机UUID
func generateUUID() string {
	// 生成随机字节序列
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		panic(err)
	}

	// 设置UUID版本和变体
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0xbf) | 0x80 // Variant 1

	// 将UUID转换为字符串并返回
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

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
	"strconv"
	"time"
)

type ClusterManager struct {
	RemoteService
	ContainerService
}

func (c *ClusterManager) getClusterIdByRole(userId uint, role string) []int {
	var ids []int
	if role != "admin" {
		db3 := database.DB
		var userClusterList []model.UserCluster
		db3.Where("user_id = ?", userId).Find(&userClusterList)
		for i := range userClusterList {
			ids = append(ids, userClusterList[i].ClusterId)
		}
	}
	return ids
}

func (c *ClusterManager) QueryCluster(ctx *gin.Context, sessions *session.Manager) {

	s := sessions.Start(ctx.Writer, ctx.Request)
	role := s.Get("role")
	userId := s.Get("userId")
	log.Println("role", role, "userId", userId)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.DB
	db2 := database.DB
	if containerId, isExist := ctx.GetQuery("containerId"); isExist {
		db = db.Where("container_id = ?", containerId)
		db2 = db2.Where("container_id = ?", containerId)
	}
	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)
	clusters := make([]model.Cluster, 0)
	ids := c.getClusterIdByRole(userId.(uint), role.(string))
	if role != "admin" {
		db.Where("id in ?", ids).Find(&clusters)
	} else {
		db.Find(&clusters)
	}

	var total int64
	if role != "admin" {
		db2.Where("id in ?", ids).Model(&model.Cluster{}).Count(&total)
	} else {
		db2.Model(&model.Cluster{}).Count(&total)
	}
	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	//var clusterVOList = make([]vo.ClusterVO, len(clusters))
	//var wg sync.WaitGroup
	//wg.Add(len(clusters))
	//for i, cluster := range clusters {
	//	go func(cluster model.Cluster, i int) {
	//		clusterVO := vo.ClusterVO{
	//			Name:            cluster.Name,
	//			ClusterName:     cluster.ClusterName,
	//			Description:     cluster.Description,
	//			ID:              cluster.ID,
	//			CreatedAt:       cluster.CreatedAt,
	//			UpdatedAt:       cluster.UpdatedAt,
	//			Ip:              cluster.Ip,
	//			Port:            cluster.Port,
	//			Username:        cluster.Username,
	//			ClusterPassword: cluster.Password,
	//		}
	//		clusterVO.GameArchive = c.GetRemoteGameArchive(cluster)
	//		clusterVO.Status = c.GetRemoteLevelStatus(cluster)
	//		clusterVOList[i] = clusterVO
	//		wg.Done()
	//	}(cluster, i)
	//}
	//wg.Wait()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       clusters,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
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
	if err != nil {
		log.Panicln(err)
	}

	cluster.ContainerId = containerId
	cluster.Uuid = containerId
	cluster.ClusterName = containerId
	cluster.Status = "init"
	cluster.ExpireTime = time.Now().Add(time.Duration(cluster.Day) * time.Hour * 24).Unix()

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

	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if clusterName == "" {
		log.Panicln("cluster is not allow null")
	}

	cluster := model.Cluster{}
	result := db.Where("cluster_name = ?", clusterName).Unscoped().Delete(&cluster)
	if result.Error != nil {
		log.Panicln(result.Error)
	}
	log.Println(cluster)
	err := c.DeleteContainer(clusterName)

	if err != nil {
		log.Panicln(err)
	}

	tx.Commit()
	return &cluster, nil
}

func (c *ClusterManager) FindClusterByUuid(uuid string) *model.Cluster {
	db := database.DB
	cluster := &model.Cluster{}
	db.Where("uuid=?", uuid).First(cluster)
	return cluster
}

func (c *ClusterManager) Restart(containerId string) error {
	err := c.RestartContainer(containerId)
	return err
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

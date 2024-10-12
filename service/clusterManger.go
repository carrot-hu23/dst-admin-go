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
)

type ClusterManager struct {
	RemoteService
	ContainerService
}

func (c *ClusterManager) GetCluster(clusterName string) *model.Cluster {
	db := database.DB
	var cluster model.Cluster
	db.Where("cluster_name = ?", clusterName).Find(&cluster)
	return &cluster
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
	if activate, isExist := ctx.GetQuery("activate"); isExist {
		boolValue := false
		if activate == "true" {
			boolValue = true
		}
		db = db.Where("activate = ?", boolValue)
		db2 = db2.Where("activate = ?", boolValue)
	}
	if levelNum, isExist := ctx.GetQuery("levelNum"); isExist {
		intValue, _ := strconv.Atoi(levelNum)
		db = db.Where("level_num = ?", intValue)
		db2 = db2.Where("level_num = ?", intValue)
	}
	if maxPlayers, isExist := ctx.GetQuery("maxPlayers"); isExist {
		intValue, _ := strconv.Atoi(maxPlayers)
		db = db.Where("max_players = ?", intValue)
		db2 = db2.Where("max_players = ?", intValue)
	}
	if core, isExist := ctx.GetQuery("core"); isExist {
		intValue, _ := strconv.Atoi(core)
		db = db.Where("core = ?", intValue)
		db2 = db2.Where("core = ?", intValue)
	}
	if memory, isExist := ctx.GetQuery("memory"); isExist {
		intValue, _ := strconv.Atoi(memory)
		db = db.Where("memory = ?", intValue)
		db2 = db2.Where("memory = ?", intValue)
	}

	db = db.Where("activate", true)
	db2 = db2.Where("activate", true)

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

func (c *ClusterManager) CreateCluster(cluster *model.Cluster) error {

	db := database.DB
	tx := db.Begin()

	//containerId, err := c.CreateContainer(*cluster)
	//if err != nil {
	//	tx.Rollback()
	//	return err
	//}

	// 生成uuid
	uuid := generateUUID()
	cluster.Uuid = uuid
	cluster.ClusterName = uuid
	cluster.Status = "init"
	cluster.Expired = false
	cluster.Activate = false
	// cluster.ExpireTime = time.Now().Add(time.Duration(cluster.Day) * time.Hour * 24).Unix()

	err := db.Create(&cluster).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
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
			log.Panicln(r)
		}
	}()

	if clusterName == "" {
		tx.Rollback()
		log.Panicln("cluster is not allow null")
	}

	cluster := model.Cluster{}
	db.Where("cluster_name= ?", clusterName).Find(&cluster)
	db.Where("cluster_name = ?", clusterName).Delete(&model.Cluster{})
	log.Println("正在删除cluster", cluster.ClusterName)
	err := c.DeleteContainer(cluster.ClusterName)

	if err != nil {
		tx.Rollback()
		log.Panicln(err)
	}

	// 删除绑定的关系
	db.Where("cluster_id = ?", cluster.ID).Delete(&[]model.UserCluster{})

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

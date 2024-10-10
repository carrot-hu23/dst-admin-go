package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ClusterApi struct{}

var clusterManager = service.ClusterManager{}

func (c *ClusterApi) GetClusterList(ctx *gin.Context) {
	clusterManager.QueryCluster(ctx, sessions)
}

func (c *ClusterApi) CreateCluster(ctx *gin.Context) {

	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	if clusterModel.Day == 0 {
		log.Panicln("过期时间不能为0")
	}
	if clusterModel.LevelNum == 0 {
		log.Panicln("世界层数不能为0")
	}
	if clusterModel.MaxPlayers == 0 {
		log.Panicln("玩家人数不能为0")
	}
	if clusterModel.Core == 0 {
		log.Panicln("cpu核数不能为0")
	}
	if clusterModel.Memory == 0 {
		log.Panicln("内存不能为0")
	}
	if clusterModel.Disk == 0 {
		log.Panicln("磁盘不能为0")
	}
	fmt.Printf("%v", clusterModel)

	var clusterList []model.Cluster

	// 批量创建
	quantity := clusterModel.Quantity
	for i := 0; i < quantity; i++ {

		cluster := model.Cluster{
			LevelNum:   clusterModel.LevelNum,
			MaxPlayers: clusterModel.MaxPlayers,
			Core:       clusterModel.Core,
			Disk:       clusterModel.Disk,
			Day:        clusterModel.Day,
			Name:       fmt.Sprintf("%s-%d", clusterModel.Name, i+1),
			Image:      clusterModel.Image,
		}
		// 计算端口
		portStart := getStartPort()
		portEnd := portStart
		cluster.Port = int(portStart)
		cluster.MasterPort = int(portStart + 1)
		portEnd = portEnd + 2
		// 冗余 5 个端口
		portEnd = portEnd + int64(cluster.LevelNum) + 5
		// 保存
		saveEndPort(portEnd)
		log.Println("正在创建cluster", cluster)
		clusterManager.CreateCluster(&cluster)
		clusterList = append(clusterList, cluster)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterList,
	})

}

func (c *ClusterApi) UpdateCluster(ctx *gin.Context) {
	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", clusterModel)
	log.Println("clusterModel", clusterModel)
	clusterManager.UpdateCluster(&clusterModel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) DeleteCluster(ctx *gin.Context) {

	clusterName := ctx.Query("clusterName")

	clusterModel, err := clusterManager.DeleteCluster(clusterName)
	log.Println("删除", clusterModel)
	if err != nil {
		log.Panicln("delete cluster error", err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) RestartCluster(ctx *gin.Context) {

	clusterName := ctx.Query("clusterName")

	err := clusterManager.RestartContainer(clusterName)
	log.Println("重启", clusterName)
	if err != nil {
		log.Panicln("restart cluster error", err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) GetCluster(ctx *gin.Context) {

	clusterName := ctx.Param("id")
	fmt.Printf("%s", clusterName)

	db := database.DB
	var cluster model.Cluster
	db.Where("cluster_name = ?", clusterName).Find(&cluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: cluster,
	})
}

func (c *ClusterApi) UpdateClusterContainer(ctx *gin.Context) {
	var payload struct {
		ClusterName string `json:"ClusterName"`
		Day         int64  `json:"day"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}

	db := database.DB
	var cluster model.Cluster
	db.Where("cluster_name = ?", payload.ClusterName).Find(&cluster)

	cluster.Day = cluster.Day + payload.Day
	cluster.ExpireTime = cluster.ExpireTime + payload.Day*24*60*60

	db.Save(&cluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: cluster,
	})

}

func getStartPort() int64 {
	version, err := fileUtils.ReadFile("./startPort")
	if err != nil {
		log.Println(err)
		return 20000
	}
	version = strings.Replace(version, "\n", "", -1)
	l, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		log.Println(err)
		return 20000
	}
	return l
}

func saveEndPort(portEnd int64) {
	fileUtils.WriterTXT(strconv.Itoa(int(portEnd)), "./startPort")
}

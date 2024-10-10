package api

import (
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	clusterManager.CreateCluster(&clusterModel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
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

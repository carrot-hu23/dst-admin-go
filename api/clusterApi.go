package api

import (
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"dst-admin-go/vo/cluster"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClusterApi struct{}

var clusterService = service.ClusterService{}
var clusterManager = service.ClusterManager{}

func (c *ClusterApi) GetGameConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterService.GetGameConfig(ctx),
	})
}

func (c *ClusterApi) SaveGameConfig(ctx *gin.Context) {

	gameConfig := cluster.GameConfig{}
	ctx.ShouldBind(&gameConfig)
	fmt.Printf("%v", gameConfig.Caves.ServerIni)
	clusterService.SaveGameConfig(ctx, &gameConfig)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *ClusterApi) GetClusterList(ctx *gin.Context) {
	clusterManager.QueryCluster(ctx)
}

func (c *ClusterApi) CreateCluster(ctx *gin.Context) {

	clusterModel := model.Cluster{}
	ctx.ShouldBind(&clusterModel)
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
	ctx.ShouldBind(&clusterModel)
	fmt.Printf("%v", clusterModel)
	clusterManager.UpdateCluster(&clusterModel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) DeleteCluster(ctx *gin.Context) {

	var id int

	if idParam, isExist := ctx.GetQuery("id"); isExist {
		id, _ = strconv.Atoi(idParam)
	}

	err := clusterManager.DeleteCluster(uint(id))
	if err != nil {
		log.Panicln("delete cluster error", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

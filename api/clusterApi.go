package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"dst-admin-go/vo/cluster"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClusterApi struct{}

var clusterService = service.ClusterService{}

func (c *ClusterApi) GetClusterConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterService.GetMultiLevelWorldConfig(),
	})
}

func (c *ClusterApi) SaveClusterConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *ClusterApi) GetGameConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterService.GetGameConfog(),
	})
}

func (c *ClusterApi) SaveGameConfig(ctx *gin.Context) {

	gameConfig := cluster.GameConfig{}
	ctx.ShouldBind(&gameConfig)
	fmt.Printf("%v", gameConfig.Caves.ServerIni)
	clusterService.SaveGameConfig(&gameConfig)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

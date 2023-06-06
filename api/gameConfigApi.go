package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameConfigApi struct {
}

var gameConfigService = service.GameConfigService{}

func (g *GameConfigApi) GetConfig(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameConfigService.GetConfig(clusterName),
	})
}

func (g *GameConfigApi) SaveConfig(ctx *gin.Context) {

	gameConfig := vo.NewGameConfigVO()
	ctx.ShouldBind(gameConfig)
	log.Println(gameConfig)
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	gameConfigService.SaveConfig(clusterName, *gameConfig)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "save dst server config success",
	})
}

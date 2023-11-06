package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GameApi struct {
}

var gameService = service.GameService{}

func (g *GameApi) UpdateGame(ctx *gin.Context) {

	log.Println("正在更新游戏。。。。。。")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	err := gameService.UpdateGame(clusterName)
	if err != nil {
		log.Panicln("更新游戏失败: ", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "update dst success",
		Data: nil,
	})
}

func (g *GameApi) GetSystemInfo(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameService.GetSystemInfo(clusterName),
	})
}

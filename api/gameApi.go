package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SpecifiedGameApi struct {
}

var gameService = service.GameService{}

func (s *SpecifiedGameApi) UpdateGame(ctx *gin.Context) {

	log.Println("正在更新游戏。。。。。。")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	gameService.UpdateGame(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "update dst success",
		Data: nil,
	})
}

func (s *SpecifiedGameApi) StartSpecifiedGame(ctx *gin.Context) {

	opType, _ := strconv.Atoi(ctx.DefaultQuery("type", "0"))

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	log.Println("正在启动指定游戏服务 type:", clusterName, opType)
	gameService.StartSpecifiedGame(clusterName, opType)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start " + clusterName + " game success",
		Data: nil,
	})
}

func (s *SpecifiedGameApi) StopSpecifiedGame(ctx *gin.Context) {

	opType, _ := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	log.Println("正在停止指定游戏服务 type:", clusterName, opType)

	gameService.StopSpecifiedGame(clusterName, opType)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "stop " + clusterName + " game success",
		Data: nil,
	})
}

func (s *SpecifiedGameApi) GetSpecifiedDashboardInfo(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameService.GetSpecifiedClusterDashboard(clusterName),
	})
}

func CreateNewClusterHome() {

}

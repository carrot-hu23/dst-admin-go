package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"dst-admin-go/vo/level"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GameLevelApi struct {
}

var gameLevelService = service.GameLevelService{}

func (g *GameLevelApi) GetLevelList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameLevelService.GetLevelsConfig(clusterName),
	})
}

func (g *GameLevelApi) DeleteLevel(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")

	gameLevelService.DeleteLevelConfig(clusterName, levelName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevelApi) CreateNewLevel(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	newLevel := service.LevelConfig{}
	err := ctx.ShouldBind(&newLevel)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	gameLevelService.CreateNewLevel(clusterName, newLevel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevelApi) GetLeveldataoverride(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameLevelService.GetLeveldataoverride(clusterName, levelName),
	})
}

func (g *GameLevelApi) GetModoverrides(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameLevelService.GetModoverrides(clusterName, levelName),
	})
}

func (g *GameLevelApi) GetServerIni(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameLevelService.GetServerIni(clusterName, levelName),
	})
}

func (g *GameLevelApi) SaveLeveldataoverride(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	var payload struct {
		LevelName         string `json:"levelName"`
		Leveldataoverride string `json:"leveldataoverride"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	gameLevelService.SaveLeveldataoverride(clusterName, payload.LevelName, payload.Leveldataoverride)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevelApi) SaveModoverrides(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	var payload struct {
		LevelName    string `json:"levelName"`
		Modoverrides string `json:"modoverrides"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	gameLevelService.SaveModoverrides(clusterName, payload.LevelName, payload.Modoverrides)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevelApi) SaveServerIni(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	var payload struct {
		LevelName string           `json:"levelName"`
		ServerIni *level.ServerIni `json:"serverIni"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	gameLevelService.SaveServerIni(clusterName, payload.LevelName, payload.ServerIni)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

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

type GameLevel2Api struct {
}

var gameLevel2Service = service.GameLevel2Service{}

func (g *GameLevel2Api) GetLevelList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameLevel2Service.GetLevelList(clusterName),
	})
}

func (g *GameLevel2Api) UpdateLevelsList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	var body struct {
		Levels []level.World `json:"levels"`
	}
	err := ctx.ShouldBind(&body)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	err = gameLevel2Service.UpdateLevels(clusterName, body.Levels)
	if err != nil {
		log.Panicln("更新世界配置失败", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevel2Api) DeleteLevel(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")

	err := gameLevel2Service.DeleteLevel(clusterName, levelName)
	if err != nil {
		log.Panicln("删除世界失败", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevel2Api) CreateNewLevel(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	newLevel := &level.World{}
	err := ctx.ShouldBind(newLevel)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	err = gameLevel2Service.CreateLevel(clusterName, newLevel)
	if err != nil {
		log.Panicln("创建世界失败", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: newLevel,
	})
}

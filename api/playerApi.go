package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayerApi struct {
}

var playerService = service.PlayerService{}

func (p *PlayerApi) GetDstPlayerList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerService.GetPlayerList(clusterName, "Master"),
	})
}

func (p *PlayerApi) GetDstAdminList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerService.GetDstAdminList(clusterName),
	})
}

func (p *PlayerApi) GetDstBlcaklistPlayerList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerService.GetDstBlacklistPlayerList(clusterName),
	})
}

func (p *PlayerApi) SaveDstAdminList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	adminListVO := vo.NewAdminListVO()
	ctx.BindJSON(adminListVO)
	playerService.SaveDstAdminList(clusterName, adminListVO.AdminList)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

func (p *PlayerApi) SaveDstBlacklistPlayerList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	blacklistVO := vo.NewBlacklistVO()
	ctx.BindJSON(blacklistVO)
	playerService.SaveDstBlacklistPlayerList(clusterName, blacklistVO.Blacklist)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

func (p *PlayerApi) DeleteDstBlacklistPlayerList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	blacklistVO := vo.NewBlacklistVO()
	ctx.BindJSON(blacklistVO)
	playerService.DeleteDstBlacklistPlayerList(clusterName, blacklistVO.Blacklist)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

func (p *PlayerApi) DeleteDstAdminListPlayerList(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	adminlistVO := vo.NewAdminListVO()
	ctx.BindJSON(adminlistVO)
	playerService.DeleteDstAdminListPlayerList(clusterName, adminlistVO.AdminList)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayerApi struct {
}

var playerService = service.PlayerService{}

func (p *PlayerApi) GetDstPlayerList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerService.GetPlayerList(),
	})
}

func (p *PlayerApi) GetDstAdminList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerService.GetDstAdminList(),
	})
}

func (p *PlayerApi) GetDstBlcaklistPlayerList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerService.GetDstBlcaklistPlayerList(),
	})
}

func (p *PlayerApi) SaveDstAdminList(ctx *gin.Context) {

	adminListVO := vo.NewAdminListVO()
	ctx.BindJSON(adminListVO)
	playerService.SaveDstAdminList(adminListVO.AdminList)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

func (p *PlayerApi) SaveDstBlacklistPlayerList(ctx *gin.Context) {

	blacklistVO := vo.NewBlacklistVO()
	ctx.BindJSON(blacklistVO)
	playerService.SaveDstBlacklistPlayerList(blacklistVO.Blacklist)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

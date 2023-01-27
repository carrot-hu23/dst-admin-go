package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDstPlayerList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func GetDstAdminList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetDstAdminList(),
	})
}

func GetDstBlcaklistPlayerList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetDstBlcaklistPlayerList(),
	})
}

func SaveDstAdminList(ctx *gin.Context) {

	adminListVO := vo.NewAdminListVO()
	ctx.BindJSON(adminListVO)
	service.SaveDstAdminList(adminListVO.AdminList)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

func SaveDstBlacklistPlayerList(ctx *gin.Context) {

	blacklistVO := vo.NewBlacklistVO()
	ctx.BindJSON(blacklistVO)
	service.SaveDstBlacklistPlayerList(blacklistVO.Blacklist)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
	})
}

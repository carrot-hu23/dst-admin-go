package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initThirdPartyRouter(router *gin.RouterGroup) {

	thirdPartyApi := api.ThirdPartyApi{}
	//第三方api转发
	router.GET("/api/dst/version", thirdPartyApi.GetDstVersion)
	router.POST("/api/dst/home/server", thirdPartyApi.GetDstHomeServerList)
	router.POST("/api/dst/home/server/detail", thirdPartyApi.GetDstHomeDetailList)
	router.GET("/api/dst/lobby/server/detail", thirdPartyApi.QueryLobbyServerDetail)
}

package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initDstConfigRouter(router *gin.RouterGroup) {

	dstConfigApi := api.DstConfigApi{}
	dstConfig := router.Group("/api/dst/config")
	{
		dstConfig.GET("", dstConfigApi.GetDstConfig)
		dstConfig.POST("", dstConfigApi.SaveDstConfig)
	}
}

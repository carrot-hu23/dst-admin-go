package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initWebLinkRouter(router *gin.RouterGroup) {

	api := api.WebLinkApi{}
	group := router.Group("/api/web/link")
	{
		group.GET("", api.GetWebLinkList)
		group.POST("", api.AddWebLink)
		group.DELETE("", api.DeleteWebLink)
	}

}

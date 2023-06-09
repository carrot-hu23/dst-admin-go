package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initSteamRouter(router *gin.RouterGroup) {

	steamApi := api.SteamApi{}
	steam := router.Group("/steam")
	{
		steam.GET("/dst/news", steamApi.DstNews)
	}

}

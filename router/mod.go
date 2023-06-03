package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initModRouter(router *gin.RouterGroup) {

	modApi := api.ModApi{}
	mod := router.Group("/api/mod")
	{
		mod.GET("/search", modApi.SearchModList)
		mod.GET("/:modId", modApi.GetModInfo)
		mod.GET("", modApi.GetMyModList)
		mod.DELETE("/:modId", modApi.DeleteMod)
	}

}

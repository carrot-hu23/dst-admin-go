package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initPlayerRouter(router *gin.RouterGroup) {

	playerApi := api.PlayerApi{}
	player := router.Group("/api/game/player")
	{
		player.GET("", playerApi.GetDstPlayerList)
		player.GET("/admin", playerApi.GetDstAdminList)
		player.POST("/admin", playerApi.SaveDstAdminList)
		player.GET("/blacklist", playerApi.GetDstBlcaklistPlayerList)
		player.POST("/blacklist", playerApi.SaveDstBlacklistPlayerList)
	}

	playerLogApi := api.PlayerLogApi{}
	playerLog := router.Group("/api/player")
	{
		playerLog.GET("/log", playerLogApi.PlayerLogQueryPage)
	}

}

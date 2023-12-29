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
		player.GET("/adminlist", playerApi.GetDstAdminList)
		player.POST("/adminlist", playerApi.SaveDstAdminList)
		player.GET("/blacklist", playerApi.GetDstBlcaklistPlayerList)
		player.POST("/blacklist", playerApi.SaveDstBlacklistPlayerList)
		player.DELETE("/blacklist", playerApi.DeleteDstBlacklistPlayerList)
		player.DELETE("/adminlist", playerApi.DeleteDstAdminListPlayerList)
	}

	playerLogApi := api.PlayerLogApi{}
	playerLog := router.Group("/api/player")
	{
		playerLog.GET("/log", playerLogApi.PlayerLogQueryPage)
		playerLog.POST("/log/delete", playerLogApi.DeletePlayerLog)

	}

}

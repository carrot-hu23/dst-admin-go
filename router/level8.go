package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func init8Level(router *gin.RouterGroup) {

	game8LevelApi := api.Game8LevelApi{}
	group := router.Group("/api/game/8level")
	{
		group.GET("/status", game8LevelApi.GetStatus)
		group.GET("/start", game8LevelApi.Start)
		group.GET("/stop", game8LevelApi.Stop)

		group.GET("/clusterIni", game8LevelApi.GetClusterIni)
		group.POST("/clusterIni", game8LevelApi.SaveClusterIni)

		group.GET("/players", game8LevelApi.GetOnlinePlayers)
		group.GET("/adminilist", game8LevelApi.GetAdministrators)
		group.GET("/whitelist", game8LevelApi.GetWhitelist)
		group.GET("/blacklist", game8LevelApi.GetBlacklist)

		group.POST("/adminilist", game8LevelApi.SaveAdminlist)
		group.POST("/whitelist", game8LevelApi.SaveWhitelist)
		group.POST("/blacklist", game8LevelApi.SaveBlacklist)

		group.GET("/command", game8LevelApi.SendCommand)

	}

}

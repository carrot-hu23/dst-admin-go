package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initLevel2(router *gin.RouterGroup) {

	gameLevelApi := api.GameLevel2Api{}
	group := router.Group("/api/cluster/level")
	{
		group.GET("", gameLevelApi.GetLevelList)
		group.PUT("", gameLevelApi.SaveLevelsList)
		group.POST("", gameLevelApi.CreateNewLevel)
		group.DELETE("", gameLevelApi.DeleteLevel)
	}

	group2 := router.Group("/api/game/8level")
	{
		group2.GET("/status", gameLevelApi.GetStatus)
		group2.GET("/start", gameLevelApi.Start)
		group2.GET("/stop", gameLevelApi.Stop)
		group2.GET("/start/all", gameLevelApi.StartAll)
		group2.GET("/stop/all", gameLevelApi.StopAll)

		group2.GET("/clusterIni", gameLevelApi.GetClusterIni)
		group2.POST("/clusterIni", gameLevelApi.SaveClusterIni)

		group2.GET("/players", gameLevelApi.GetOnlinePlayers)
		group2.GET("/players/all", gameLevelApi.GetAllOnlinePlayers)

		group2.GET("/adminilist", gameLevelApi.GetAdministrators)
		group2.GET("/whitelist", gameLevelApi.GetWhitelist)
		group2.GET("/blacklist", gameLevelApi.GetBlacklist)

		group2.POST("/adminilist", gameLevelApi.SaveAdminlist)
		group2.POST("/whitelist", gameLevelApi.SaveWhitelist)
		group2.POST("/blacklist", gameLevelApi.SaveBlacklist)

		group2.GET("/command", gameLevelApi.SendCommand)

		group2.GET("/udp/port", gameLevelApi.GetScanUDPPorts)
	}

	preinstallApi := api.PreinstallApi{}
	group3 := router.Group("/api/game/preinstall")
	{
		group3.GET("", preinstallApi.UsePreinstall)
	}

	shareApi := api.ShareApi{}
	group4 := router.Group("/api/share")
	{
		group4.GET("/keyCer", shareApi.GetKeyCerApi)
		group4.GET("/keyCer/reflush", shareApi.ReflushKeyCerApi)
		group4.GET("/keyCer/enable", shareApi.EnableKeyCerApi)

		group4.POST("/cluster/import", shareApi.ImportClusterConfig)
	}

	group5 := router.Group("/share")
	{
		group5.GET("/cluster", shareApi.ShareClusterConfig)

	}

}

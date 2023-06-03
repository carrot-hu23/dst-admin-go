package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func InitGameRouter(router *gin.RouterGroup) {

	gameApi := api.GameConsoleApi{}
	game := router.Group("/api/game")
	{
		game.GET("/update", gameApi.UpdateGame)

		game.GET("/sent/broadcast", gameApi.SentBroadcast)
		game.GET("/kick/player", gameApi.KickPlayer)
		game.GET("/kill/player", gameApi.KillPlayer)
		game.GET("/respawn/player", gameApi.RespawnPlayer)
		game.GET("/rollback", gameApi.RollBack)
		game.GET("/regenerateworld", gameApi.Regenerateworld)
		game.POST("/master/console", gameApi.MasterConsole)
		game.POST("/caves/console", gameApi.CavesConsole)
		game.GET("/operate/player", gameApi.OperatePlayer)
		game.GET("/backup/restore", gameApi.RestoreBackup)

		game.GET("/archive", gameApi.GetGameArchive)
		game.GET("/clean", gameApi.CleanWorld)
	}

	gameConfigApi := api.GameConfigApi{}
	gameConfig := router.Group("/api/game/config")
	{
		gameConfig.GET("", gameConfigApi.GetConfig)
		gameConfig.POST("", gameConfigApi.SaveConfig)
	}

	specifiedGameApi := api.SpecifiedGameApi{}
	specified := router.Group("/api/game/specified")
	{
		specified.GET("/dashboard", specifiedGameApi.GetSpecifiedDashboardInfo)
		specified.GET("/start", specifiedGameApi.StartSpecifiedGame)
		specified.GET("/stop", specifiedGameApi.StopSpecifiedGame)
	}

}

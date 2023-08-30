package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func InitGameRouter(router *gin.RouterGroup) {

	gameConsoleApi := api.GameConsoleApi{}
	gameConsole := router.Group("/api/game")
	{
		gameConsole.GET("/sent/broadcast", gameConsoleApi.SentBroadcast)
		gameConsole.GET("/kick/player", gameConsoleApi.KickPlayer)
		gameConsole.GET("/kill/player", gameConsoleApi.KillPlayer)
		gameConsole.GET("/respawn/player", gameConsoleApi.RespawnPlayer)
		gameConsole.GET("/rollback", gameConsoleApi.RollBack)
		gameConsole.GET("/regenerateworld", gameConsoleApi.Regenerateworld)
		gameConsole.POST("/master/console", gameConsoleApi.MasterConsole)
		gameConsole.POST("/caves/console", gameConsoleApi.CavesConsole)
		gameConsole.GET("/operate/player", gameConsoleApi.OperatePlayer)
		gameConsole.GET("/backup/restore", gameConsoleApi.RestoreBackup)

		gameConsole.GET("/archive", gameConsoleApi.GetGameArchive)
		gameConsole.GET("/clean", gameConsoleApi.CleanWorld)
		gameConsole.GET("/clean/level", gameConsoleApi.CleanLevel)
		gameConsole.GET("/announce/setting", gameConsoleApi.GetAnnounceSetting)
		gameConsole.POST("/announce/setting", gameConsoleApi.SaveAnnounceSetting)
		gameConsole.GET("/level/server/log", gameConsoleApi.ReadLevelServeLog)
		gameConsole.GET("/level/server/chat/log", gameConsoleApi.ReadLevelServeChatLog)
	}

	gameConfigApi := api.GameConfigApi{}
	gameConfig := router.Group("/api/game/config")
	{
		gameConfig.GET("", gameConfigApi.GetConfig)
		gameConfig.POST("", gameConfigApi.SaveConfig)
	}

	gameApi := api.GameApi{}
	game := router.Group("/api/game")
	{
		game.GET("/update", gameApi.UpdateGame)
		game.GET("/dashboard", gameApi.GetDashboardInfo)
		game.GET("/start", gameApi.StartGame)
		game.GET("/stop", gameApi.StopGame)
	}

}

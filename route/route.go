package route

import (
	"dst-admin-go/api"
	"dst-admin-go/handler"

	"github.com/gin-gonic/gin"
)

func loadStaticeFile(app *gin.Engine) {
	// dir, _ := os.Getwd()
	app.LoadHTMLGlob("dist/index.html") // 添加入口index.html
	//r.LoadHTMLFiles("dist//*") // 添加资源路径
	app.Static("/assets", "./dist/assets")
	app.Static("/misc", "./dist/misc")
	app.Static("/static/js", "./dist/static/js")                         // 添加资源路径
	app.Static("/static/css", "./dist/static/css")                       // 添加资源路径
	app.Static("/static/img", "./dist/static/img")                       // 添加资源路径
	app.Static("/static/fonts", "./dist/static/fonts")                   // 添加资源路径
	app.Static("/static/media", "./dist/static/media")                   // 添加资源路径
	app.StaticFile("/favicon.ico", "./dist/favicon.ico")                 // 添加资源路径
	app.StaticFile("/asset-manifest.json", "./dist/asset-manifest.json") // 添加资源路径
	app.StaticFile("/", "./dist/index.html")                             //前端接口
}

func NewRoute() *gin.Engine {

	app := gin.Default()

	app.Use(handler.Recover)
	app.Use(handler.ShellInjectionInterceptor())
	app.Use(handler.Authentucation())

	// app.Use(handler.CheckDstHandler())

	app.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello! Dont starve together")
	})

	app.POST("/api/login", api.Login)
	app.GET("/api/logout", api.Logout)
	app.POST("/api/change/password", api.ChangePassword)

	app.GET("/api/init", api.CheckIsFirst)
	app.POST("/api/init", api.InitFirst)

	app.GET("/ws", api.HandlerWS)

	user := app.Group("/api/user")
	{
		user.GET("", api.GetUserInfo)
	}

	dashboard := app.Group("/api/dashboard")
	{
		dashboard.GET("", api.GetDashboardInfo)
	}

	gameConfig := app.Group("/api/game/config")
	{
		gameConfig.GET("", api.GetConfig)
		gameConfig.POST("", api.SaveConfig)
	}

	game := app.Group("/api/game")
	{
		game.GET("/update", api.UpdateGame)
		game.GET("/start", api.StartGame)
		game.GET("/stop", api.StoptGame)
		game.GET("/start/master", api.StartMaster)
		game.GET("/start/caves", api.StartCaves)
		game.GET("/stop/master", api.StopMaster)
		game.GET("/stop/caves", api.StopCaves)

		game.GET("/sent/broadcast", api.SentBroadcast)
		game.GET("/kick/player", api.KickPlayer)
		game.GET("/kill/player", api.KillPlayer)
		game.GET("/respawn/player", api.RespawnPlayer)
		game.GET("/rollback", api.RollBack)
		game.GET("/regenerateworld", api.Regenerateworld)
		game.GET("/master/console", api.MasterConsole)
		game.GET("/caves/console", api.CavesConsole)
		game.GET("/operate/player", api.OperatePlayer)
		game.GET("/backup/restore", api.RestoreBackup)

		game.GET("/archive", api.GetGameArchive)
	}

	player := app.Group("/api/game/player")
	{
		player.GET("", api.GetDstPlayerList)
		player.GET("/admin", api.GetDstAdminList)
		player.POST("/admin", api.SaveDstAdminList)
		player.GET("/blacklist", api.GetDstBlcaklistPlayerList)
		player.POST("/blacklist", api.SaveDstBlacklistPlayerList)
	}

	dstConfig := app.Group("/api/dst/config")
	{
		dstConfig.GET("", api.GetDstConfig)
		dstConfig.POST("", api.SaveDstConfig)
	}

	backup := app.Group("/api/game/backup")
	{
		backup.GET("", api.GetBackupList)
		backup.POST("", api.CreateBackup)
		backup.DELETE("", api.DeleteBackup)
		backup.PUT("", api.RenameBackup)
		backup.GET("/download", api.DownloadBackup)
		backup.POST("/upload", api.UploadBackup)
	}

	//第三方api转发
	app.GET("/api/dst/version", api.GetDstVersion)
	app.POST("/api/dst/home/server", api.GetDstHomeServerList)
	app.POST("/api/dst/home/server/detail", api.GetDstHomeDetailList)

	playerLog := app.Group("/api/player")
	{
		playerLog.GET("/log", api.PlayerLogQueryPage)
	}

	mod := app.Group("/api/mod")
	{
		mod.GET("/search", api.SearchModList)
		mod.GET("/:modId", api.GetModInfo)
		mod.GET("", api.GetMyModList)
		mod.DELETE("/:modId", api.DeleteMod)
	}

	statistics := app.Group("/api/statistics")
	{
		statistics.GET("/active/user", api.CountActiveUser)
		statistics.GET("/top/death", api.TopDeaths)
		statistics.GET("/top/login", api.TopUserLoginimes)
		statistics.GET("/top/active", api.TopUserActiveTimes)

		statistics.GET("/rate/role", api.CountRoleRate)
	}

	app.Any("/app/:name/*path", api.NewProxy)

	proxyApp := app.Group("/api/proxy")
	{
		proxyApp.GET("", api.GetProxyEntity)
		proxyApp.POST("", api.CreateProxyEntity)
		proxyApp.PUT("", api.UpdateProxyEntity)
		proxyApp.DELETE("", api.DeleteProxyEntity)
	}
	loadStaticeFile(app)
	return app
}

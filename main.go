package main

import (
	"dst-admin-go/api"
	"dst-admin-go/handler"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const log_path = "./dst-admin-go.log"

var f *os.File
var port = ":8080"

func arg() {
	for idx, args := range os.Args {
		split := strings.Split(args, "=")
		if len(split) == 2 {
			argName := strings.TrimSpace(split[0])
			if argName == "port" {
				argValue := strings.TrimSpace(split[1])
				port = argValue
			}

		}
		fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	}
}

func logInit() {
	var err error
	f, err = os.OpenFile(log_path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	gin.ForceConsoleColor()
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	arg()
	logInit()
}

func main() {

	defer func() {
		f.Close()
	}()

	fmt.Println(":pig, 你是好人")
	app := gin.Default()

	app.Use(handler.Recover)
	app.Use(handler.Authentucation())
	// app.Use(handler.CheckDstHandler())

	app.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello! Dont starve together")
	})

	app.POST("/api/login", api.Login)
	app.GET("/api/logout", api.Logout)
	app.POST("/api/change/password", api.ChangePassword)

	app.GET("/first", api.CheckIsFirst)
	app.POST("/init", api.InitFirst)

	app.GET("/ws", api.HandlerWS)

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
		game.POST("/update", api.UpdateGame)
		game.GET("/start", api.StartGame)
		game.GET("/stop", api.StoptGame)
		game.GET("/start/master", api.StartMaster)
		game.GET("/start/caves", api.StartCaves)
		game.GET("/stop/master", api.StopMaster)
		game.GET("/stop/caves", api.StopCaves)

		game.GET("/sent/broadcast", api.SentBroadcast)
		game.GET("/kick/player", api.KickPlayer)
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

	// dir, _ := os.Getwd()
	app.LoadHTMLGlob("dist/index.html") // 添加入口index.html
	//r.LoadHTMLFiles("dist//*") // 添加资源路径
	app.Static("/static/js", "./dist/static/js")                         // 添加资源路径
	app.Static("/static/css", "./dist/static/css")                       // 添加资源路径
	app.Static("/static/img", "./dist/static/img")                       // 添加资源路径
	app.Static("/static/fonts", "./dist/static/fonts")                   // 添加资源路径
	app.Static("/static/media", "./dist/static/media")                   // 添加资源路径
	app.StaticFile("/favicon.ico", "./dist/favicon.ico")                 // 添加资源路径
	app.StaticFile("/asset-manifest.json", "./dist/asset-manifest.json") // 添加资源路径
	app.StaticFile("/", "./dist/index.html")                             //前端接口

	//第三方api转发
	app.GET("/api/dst/version", api.GetDstVersion)
	app.POST("/api/dst/home/server", api.GetDstHomeServerList)
	app.POST("/api/dst/home/server/detail", api.GetDstHomeDetailList)

	app.Run(port)
}

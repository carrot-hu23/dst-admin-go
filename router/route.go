package router

import (
	"dst-admin-go/middleware"

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

	app.Use(middleware.Recover)
	app.Use(middleware.ShellInjectionInterceptor())
	app.Use(middleware.Authentucation())

	// app.Use(middleware.CheckDstHandler())

	app.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello! Dont starve together")
	})
	router := app.Group("")
	initBackupRouter(router)
	initClusterRouter(router)
	initDstConfigRouter(router)
	initFirstRouter(router)
	InitGameRouter(router)
	initLoginRouter(router)
	initModRouter(router)
	initPlayerRouter(router)
	initProxyRouter(router)
	initStatisticsRouter(router)
	initThirdPartyRouter(router)
	initWsRouter(router)

	initStaticFile(app)

	return app
}

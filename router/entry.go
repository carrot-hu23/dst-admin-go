package router

import (
	"dst-admin-go/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"time"
)

func NewRoute() *gin.Engine {

	app := gin.Default()

	store := memstore.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(60 * 24 * 7 * time.Minute.Seconds()),
		HttpOnly: true,
	})
	app.Use(sessions.Sessions("token", store))

	app.Use(middleware.Recover)
	// app.Use(middleware.ShellInjectionInterceptor())
	app.Use(middleware.Authentication())

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
	initStatisticsRouter(router)
	initThirdPartyRouter(router)
	initWsRouter(router)
	initSteamRouter(router)
	initTimedTaskRouter(router)

	initAutoCheck(router)

	initWebLinkRouter(router)
	initWebhookRouter(router)

	initLevel2(router)

	initFile(router)
	initDstGenMapRouter(router)

	initStaticFile(app)

	return app
}

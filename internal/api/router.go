package api

import (
	"dst-admin-go/internal/api/handler"
	"dst-admin-go/internal/collect"
	"dst-admin-go/internal/config"
	"dst-admin-go/internal/middleware"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/backup"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/dstMap"
	"dst-admin-go/internal/service/game"
	"dst-admin-go/internal/service/gameArchive"
	"dst-admin-go/internal/service/gameConfig"
	"dst-admin-go/internal/service/level"
	"dst-admin-go/internal/service/levelConfig"
	"dst-admin-go/internal/service/login"
	"dst-admin-go/internal/service/mod"
	"dst-admin-go/internal/service/player"
	"dst-admin-go/internal/service/update"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"gorm.io/gorm"
)

func NewRoute(cfg *config.Config, db *gorm.DB) *gin.Engine {
	app := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   int(60 * 24 * 7 * time.Minute.Seconds()),
		HttpOnly: true,
	})
	app.Use(sessions.Sessions("token", store))
	app.Use(middleware.Recover)

	app.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello! Dont starve together")
	})
	// Swagger UI
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	RegisterStaticFile(app)
	Register(cfg, db, app.Group(""))
	return app
}

func initCollectors(archive *archive.PathResolver, dstConfigService dstConfig.Config) {
	getDstConfig, err := dstConfigService.GetDstConfig("MyDediServer")
	if err != nil {
		return
	}
	clusterName := getDstConfig.Cluster
	newCollect := collect.NewCollect(archive.ClusterPath(clusterName), clusterName)
	collect.Collector = newCollect
	collect.Collector.StartCollect()
}

func RegisterStaticFile(app *gin.Engine) {

	defer func() {
		if r := recover(); r != nil {
		}
	}()
	app.Use(func(context *gin.Context) {
		context.Writer.Header().Set("Cache-Control", "public, max-age=30672000")
	})
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
	app.StaticFile("/", "./dist/index.html")
}

func Register(cfg *config.Config, db *gorm.DB, router *gin.RouterGroup) {

	// service
	dstConfigService := dstConfig.NewDstConfig(db)
	updateService := update.NewUpdateService(dstConfigService)
	resolverService, _ := archive.NewPathResolver(dstConfigService)
	loginService := login.NewLoginService(cfg)
	levelConfigUtils := levelConfig.NewLevelConfigUtils(resolverService)
	gameProcess := game.NewGame(dstConfigService, levelConfigUtils)

	gameConfigService := gameConfig.NewGameConfig(resolverService, levelConfigUtils)
	backupService := backup.NewBackupService(resolverService, dstConfigService, gameProcess)
	levelService := level.NewLevelService(gameProcess, dstConfigService, resolverService, levelConfigUtils)
	playerService := player.NewPlayerService(resolverService)
	gameArchiveService := gameArchive.NewGameArchive(gameConfigService, levelService, resolverService)
	modService := mod.NewModService(db, dstConfigService, resolverService)

	dstMapGenerator := dstMap.NewDSTMapGenerator()

	// init
	initCollectors(resolverService, dstConfigService)

	//  handler
	updateHandler := handler.NewUpdateHandler(updateService)
	gameHandler := handler.NewGameHandler(gameProcess, levelService, gameArchiveService, levelConfigUtils, resolverService)
	gameConfigHandler := handler.NewGameConfigHandler(gameConfigService)
	dstConfigHandler := handler.NewDstConfigHandler(dstConfigService, resolverService)
	loginHandler := handler.NewLoginHandler(loginService)
	backupHandler := handler.NewBackupHandler(backupService)
	levelHandler := handler.NewLevelHandler(levelService)
	playerHandler := handler.NewPlayerHandler(playerService, gameProcess)
	levelLogHandler := handler.NewLevelLogHandler(resolverService)
	kvHandler := handler.NewKvHandler(db)
	dstApiHandler := handler.NewDstApiHandler()
	dstMapHandler := handler.NewDstMapHandler(resolverService, dstMapGenerator)
	playerLogHandler := handler.NewPlayerLogHandler()
	statisticsHandler := handler.NewStatisticsHandler()
	modHandler := handler.NewModHandler(modService, dstConfigService)

	// 中间件
	router.Use(middleware.Authentication(loginService))
	router.Use(middleware.ClusterMiddleware(dstConfigService))

	//  route
	updateHandler.RegisterRoute(router)
	gameHandler.RegisterRoute(router)
	gameConfigHandler.RegisterRoute(router)
	dstConfigHandler.RegisterRoute(router)
	loginHandler.RegisterRoute(router)
	backupHandler.RegisterRoute(router)
	levelHandler.RegisterRoute(router)
	playerHandler.RegisterRoute(router)
	levelLogHandler.RegisterRoute(router)
	kvHandler.RegisterRoute(router)
	dstApiHandler.RegisterRoute(router)
	dstMapHandler.RegisterRoute(router)
	playerLogHandler.RegisterRoute(router)
	statisticsHandler.RegisterRoute(router)
	modHandler.RegisterRoute(router)

}

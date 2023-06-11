package main

import (
	"dst-admin-go/config/global"
	"dst-admin-go/initConfig"
	"dst-admin-go/router"
	"github.com/gin-contrib/pprof"
)

func init() {
	initConfig.Init()
}

func main() {

	app := router.NewRoute()
	pprof.Register(app)
	app.Run(":" + global.Config.Port)

}

package main

import (
	"dst-admin-go/bootstrap"
	"dst-admin-go/config/global"
	"dst-admin-go/router"
	"fmt"
	"github.com/gin-contrib/gzip"
)

func main() {
	bootstrap.Init()
	app := router.NewRoute()
	app.Use(gzip.Gzip(gzip.BestCompression))
	err := app.Run(":" + global.Config.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}

}

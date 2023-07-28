package main

import (
	"dst-admin-go/bootstrap"
	"dst-admin-go/config/global"
	"dst-admin-go/router"
	"fmt"
)

func main() {
	bootstrap.Init()
	app := router.NewRoute()
	err := app.Run(":" + global.Config.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}

}

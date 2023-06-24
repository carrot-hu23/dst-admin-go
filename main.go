package main

import (
	"dst-admin-go/config/global"
	"dst-admin-go/initConfig"
	"dst-admin-go/router"
	"fmt"
)

func init() {
	initConfig.Init()
}

func main() {

	app := router.NewRoute()
	err := app.Run(":" + global.Config.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}

}

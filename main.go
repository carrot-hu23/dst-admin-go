package main

import (
	"dst-admin-go/config/global"
	"dst-admin-go/initConfig"
	"dst-admin-go/router"
	"embed"
	"fmt"
)

func init() {
	initConfig.Init()
}

// 嵌入为一个文件系统 新的文件系统FS
//
//go:embed dist
//go:embed static
var f embed.FS

func main() {

	app := router.NewRoute()
	err := app.Run(":" + global.Config.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}

}

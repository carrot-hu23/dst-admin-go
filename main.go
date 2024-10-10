package main

import (
	"dst-admin-go/bootstrap"
	"dst-admin-go/config/global"
	"dst-admin-go/router"
	"dst-admin-go/schedule"
	"fmt"
	"time"
)

func main() {
	bootstrap.Init()

	// 创建一个 time.Ticker
	ticker := time.NewTicker(time.Duration(global.Config.Collect) * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				schedule.CollectContainerStatus()
			}
		}
	}()

	app := router.NewRoute()
	err := app.Run(":" + global.Config.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}

}

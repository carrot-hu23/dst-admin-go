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
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("定时任务执行时间:", t)
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

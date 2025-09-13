package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dst-admin-go/bootstrap"
	"dst-admin-go/config/global"
	"dst-admin-go/router"

	"github.com/gin-contrib/gzip"
)

func main() {
	bootstrap.Init()
	app := router.NewRoute()
	app.Use(gzip.Gzip(gzip.BestCompression))

	errCh := make(chan error)
	if global.Config.Port != "" {
		go func(errCh chan<- error) {
			err := app.Run(global.Config.BindAddress + ":" + global.Config.Port)
			if err != nil {
				errCh <- err
			}
		}(errCh)
	}
	if global.Config.SecurePort != "" {
		go func(errCh chan<- error) {
			err := app.RunTLS(global.Config.SecureBindAddress+":"+global.Config.SecureBindAddress,
				global.Config.TlsCertPath, global.Config.TlsKeyPath)
			if err != nil {
				errCh <- err
			}
		}(errCh)
	}
	// 系统退出指令
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 监听出错通道和系统退出通道
	select {
	case err := <-errCh:
		fmt.Println("启动失败！！！", err)
	case <-quit:
		log.Println("正在关闭服务器...")
	}
}

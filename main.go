package main

import (
	"dst-admin-go/bootstrap"
	"dst-admin-go/config/global"
	"dst-admin-go/router"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	bootstrap.Init()
	dstCli()

	app := router.NewRoute()
	err := app.Run(":" + global.Config.Port)
	if err != nil {
		fmt.Println("启动失败！！！", err)
	}

}

func dstCli() {
	if runtime.GOOS != "windows" {
		return
	}
	// 启动 dstcli.exe 进程
	cmd := exec.Command("cmd", "/C", "main.exe")
	// 创建一个缓冲区用于捕获输出
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
	log.Println("正在启动 main.exe")

	// 捕获系统信号，以便在关闭时关闭 dstcli.exe 进程
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("Received termination signal. Shutting down...")

		// 关闭 dstcli.exe 进程
		err := cmd.Process.Kill()
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()

}

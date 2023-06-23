package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant/consts"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path/filepath"
)

type SystemApi struct{}

type SystemConfig struct {
	gorm.Model
	Steamcmd string `json:"steamcmd"`
}

func (s *SystemApi) GetConfig(ctx *gin.Context) {

	db := database.DB
	config := SystemConfig{}
	db.First(&config)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: config,
	})
}

func (s *SystemApi) SaveConfig(ctx *gin.Context) {

	var body struct {
		Steamcmd string `json:"steamcmd"`
	}
	err := ctx.ShouldBind(&body)
	if err != nil {
		log.Panicln("参数错误")
	}

	db := database.DB
	config := SystemConfig{}
	db.First(&config)
	config.Steamcmd = body.Steamcmd

	if config.ID == 0 {
		c := SystemConfig{}
		c.Steamcmd = body.Steamcmd
		db.Create(&c)
	} else {
		db.Save(&config)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (s *SystemApi) InstallSteamCmd(ctx *gin.Context) {

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	ctx.Header("Access-Control-Allow-Origin", "*")

	// 使用一个channel来接收SSE事件
	eventCh := make(chan string)
	stopCh := make(chan byte)

	defer func() {
		if err := recover(); err != nil {
			log.Println("安装依赖错误:", err)
			fmt.Fprintf(ctx.Writer, "data: 安装依赖错误 \n\n")
		}
		close(eventCh)
		close(stopCh)
	}()

	// 在单独的goroutine中发送SSE事件
	go func() {
		s.handle(eventCh, stopCh)
	}()

	// 循环读取channel中的事件并发送给客户端
	for {
		select {
		case event := <-eventCh:
			_, err := fmt.Fprintf(ctx.Writer, event)
			if err != nil {
				// 处理错误情况，例如日志记录或返回错误响应
				fmt.Println("Error writing SSE event:", err)
				return
			}
			ctx.Writer.Flush()
		case <-stopCh:
			return
		case <-ctx.Writer.CloseNotify():
			// 如果客户端断开连接，则停止发送事件
			return
		}
	}

}

// 检测是否已经安装了 steamcmd
func installCmd2(eventCh chan string, stopCh chan byte) error {

	db := database.DB
	systemConfig := SystemConfig{}
	db.First(&systemConfig)
	var steamCmdPath string
	if systemConfig.Steamcmd == "" {
		eventCh <- "data: steamcmd 默认安装在 " + consts.HomePath + "。。。\n\n"
		steamCmdPath = consts.HomePath
		systemConfig.Steamcmd = filepath.Join(consts.HomePath, "steamcmd")
		db.Save(systemConfig)

	} else {
		eventCh <- "data: steamcmd 即将安装在 " + steamCmdPath + "。。。\n\n"
	}

	// 直接调用脚本安装
	scriptPath := "./static/script/install_steamcmd2.sh"
	err := shellUtils.Chmod(scriptPath)
	if err != nil {
		log.Panicln("设置steamcmd脚本权限错误", err)
	}
	err = commandShell(eventCh, scriptPath, steamCmdPath, consts.HomePath)
	if err != nil {
		eventCh <- "data: 安装steamcmd失败！！！ \n\n"
		return err
	}
	eventCh <- "data: 安装steamcmd成功！！！ \n\n"
	return nil
}

func (s *SystemApi) handle(eventCh chan string, stopCh chan byte) {
	err := installDependence(eventCh, stopCh)
	if err != nil {
		stopCh <- 1
		return
	}
	err = installCmd2(eventCh, stopCh)
	if err != nil {
		stopCh <- 1
		return
	}

	stopCh <- 1
}

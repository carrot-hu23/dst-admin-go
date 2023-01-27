package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const first = "./first"

type InitData struct {
	User      *vo.UserVO                `json:"user"`
	DstConfig *dstConfigUtils.DstConfig `json:"dstConfig"`
}

func InitFirst(ctx *gin.Context) {

	exist := fileUtils.Exists(first)
	if exist {
		log.Panicln("非法请求")
	}

	initData := &InitData{}
	ctx.Bind(initData)

	username := initData.User.Username
	password := initData.User.Password
	service.ChangeUser(username, password)

	dstConfig := dstConfigUtils.DstConfig{
		Steamcmd:            initData.DstConfig.Steamcmd,
		Force_install_dir:   initData.DstConfig.Force_install_dir,
		DoNotStarveTogether: initData.DstConfig.DoNotStarveTogether,
		Cluster:             initData.DstConfig.Cluster,
	}

	dstConfigUtils.SaveDstConfig(&dstConfig)

	fileUtils.CreateFile(first)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func CheckIsFirst(ctx *gin.Context) {

	exist := fileUtils.Exists(first)

	code := 200
	msg := "is first"
	if exist {
		code = 400
		msg = "is not first"
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

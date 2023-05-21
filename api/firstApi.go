package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const first = "./first"

func InitFirst(ctx *gin.Context) {

	exist := fileUtils.Exists(first)
	if exist {
		log.Panicln("非法请求")
	}

	initData := &service.InitDstData{}
	ctx.Bind(initData)

	service.InitDstEnv(initData, ctx)

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

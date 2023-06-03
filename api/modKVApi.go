package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
*
获取用户mod设置偏好

目前不做这个
*/
func GetUserModKV(ctx *gin.Context) {

	userId := ctx.Query("userId")
	log.Println("userId", userId)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetModKV(),
	})
}

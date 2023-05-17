package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClusterConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetmultiLevelWorldConfig(),
	})
}

func SaveClusterConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

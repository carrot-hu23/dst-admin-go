package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"dst-admin-go/vo/cluster"
	"log"
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

func GetGameConfog(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetGameConfog(),
	})
}

func SaveGameConfog(ctx *gin.Context) {

	gameConfig := cluster.GameConfig{}
	ctx.ShouldBind(&gameConfig)
	log.Println(gameConfig)
	service.SaveGameConfig(&gameConfig)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

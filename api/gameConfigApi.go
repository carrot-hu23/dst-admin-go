package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: service.GetConfig(),
	})
}

func SaveConfig(ctx *gin.Context) {

	gameConfig := vo.NewGameConfigVO()
	ctx.ShouldBind(gameConfig)
	log.Println(gameConfig)
	service.SaveConfig(*gameConfig)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "save dst server config success",
	})
}

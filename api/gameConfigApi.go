package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameConfigApi struct {
}

var gameConfigService = service.GameConfigService{}

func (g *GameConfigApi) GetConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameConfigService.GetConfig(ctx),
	})
}

func (g *GameConfigApi) SaveConfig(ctx *gin.Context) {

	gameConfig := vo.NewGameConfigVO()
	ctx.ShouldBind(gameConfig)
	log.Println(gameConfig)
	gameConfigService.SaveConfig(ctx, *gameConfig)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "save dst server config success",
	})
}

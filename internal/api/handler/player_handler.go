package handler

import (
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/service/game"
	"dst-admin-go/internal/service/player"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService *player.PlayerService
	gameProcess   game.Process
}

func NewPlayerHandler(playerService *player.PlayerService, gameProcess game.Process) *PlayerHandler {
	return &PlayerHandler{
		playerService: playerService,
		gameProcess:   gameProcess,
	}
}

func (p *PlayerHandler) GetPlayerList(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": p.playerService.GetPlayerList(clusterName, "Master", p.gameProcess),
	})
}

func (p *PlayerHandler) GetPlayerAllList(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": p.playerService.GetPlayerAllList(clusterName, p.gameProcess),
	})
}

func (p *PlayerHandler) RegisterRoute(router *gin.RouterGroup) {
	player := router.Group("/api/game/8level/players")
	{
		player.GET("", p.GetPlayerList)
		player.GET("/all", p.GetPlayerAllList)
	}
}

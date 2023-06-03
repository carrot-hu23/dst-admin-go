package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initWsRouter(router *gin.RouterGroup) {
	wsApi := api.WebSocketApi{}
	router.GET("/ws", wsApi.HandlerWS)
}

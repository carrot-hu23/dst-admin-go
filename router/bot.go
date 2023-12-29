package router

import (
	"dst-admin-go/bot/kook"
	"github.com/gin-gonic/gin"
)

func initBotRouter(router *gin.RouterGroup) {

	kookBotApi := kook.KookBotApi{}
	kookBot := router.Group("/bot/webhook/kook")
	{
		kookBot.POST("", kookBotApi.AuthKookWebHook)
	}

}

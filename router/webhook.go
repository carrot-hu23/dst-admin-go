package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initWebhookRouter(router *gin.RouterGroup) {

	webhookApi := api.WebhookApi{}
	//第三方api转发
	router.POST("/webhook", webhookApi.Webhook)
}

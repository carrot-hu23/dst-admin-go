package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func InitLogRouter(router *gin.RouterGroup) {

	logApi := api.LogApi{}
	group := router.Group("/api/game/log")
	{
		group.GET("/stream", logApi.Stream)
	}

}

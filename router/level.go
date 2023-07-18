package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initLevel(router *gin.RouterGroup) {

	gameLevelApi := api.GameLevelApi{}
	group := router.Group("/api/game/level")
	{
		group.GET("", gameLevelApi.GetLevelList)
		group.POST("", gameLevelApi.CreateNewLevel)
		group.DELETE("", gameLevelApi.DeleteLevel)

		group.GET("/leveldataoverride", gameLevelApi.GetLeveldataoverride)
		group.GET("/modoverrides", gameLevelApi.GetModoverrides)
		group.GET("/serverIni", gameLevelApi.GetServerIni)

		group.POST("/leveldataoverride", gameLevelApi.SaveLeveldataoverride)
		group.POST("/modoverrides", gameLevelApi.SaveModoverrides)
		group.POST("/serverIni", gameLevelApi.SaveServerIni)
	}

}

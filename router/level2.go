package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initLevel2(router *gin.RouterGroup) {

	gameLevelApi := api.GameLevel2Api{}
	group := router.Group("/api/cluster/level")
	{
		group.GET("", gameLevelApi.GetLevelList)
		group.PUT("", gameLevelApi.UpdateLevelsList)
		group.POST("", gameLevelApi.CreateNewLevel)
		group.DELETE("", gameLevelApi.DeleteLevel)
	}

}

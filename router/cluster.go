package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initClusterRouter(router *gin.RouterGroup) {

	clusterApi := api.ClusterApi{}
	clusterConfig := router.Group("/api/cluster/config")
	{
		clusterConfig.GET("", clusterApi.GetClusterConfig)
		clusterConfig.POST("", clusterApi.SaveClusterConfig)
	}

	clusterGameConfig := router.Group("/api/cluster/game/config")
	{
		clusterGameConfig.GET("", clusterApi.GetGameConfig)
		clusterGameConfig.POST("", clusterApi.SaveGameConfig)
	}

}

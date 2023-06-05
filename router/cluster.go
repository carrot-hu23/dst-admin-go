package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initClusterRouter(router *gin.RouterGroup) {

	clusterApi := api.ClusterApi{}
	cluster := router.Group("/api/cluster")
	{
		cluster.GET("", clusterApi.GetClusterList)
		cluster.POST("", clusterApi.CreateCluster)
		cluster.PUT("", clusterApi.UpdateCluster)
		cluster.DELETE("", clusterApi.DeleteCluster)
	}

	clusterGameConfig := router.Group("/api/cluster/game/config")
	{
		clusterGameConfig.GET("", clusterApi.GetGameConfig)
		clusterGameConfig.POST("", clusterApi.SaveGameConfig)
	}

}

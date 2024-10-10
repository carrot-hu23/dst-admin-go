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
		cluster.GET("/:id", clusterApi.GetCluster)
		cluster.GET("/restart", clusterApi.RestartCluster)
		cluster.POST("", clusterApi.CreateCluster)
		cluster.PUT("", clusterApi.UpdateCluster)
		cluster.DELETE("", clusterApi.DeleteCluster)

		cluster.PUT("/container", clusterApi.UpdateClusterContainer)
	}

	activate := router.Group("/activate")
	{
		activate.GET("", clusterApi.GetCluster)
		activate.POST("/bind", clusterApi.BindCluster)
	}
}

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
		cluster.POST("/remote", clusterApi.FetchRemoteClusterList)
	}

	kvApi := api.KvApi{}
	kv := router.Group("/api/kv")
	{
		kv.GET("", kvApi.GetKv)
		kv.POST("", kvApi.SaveKv)
	}
}

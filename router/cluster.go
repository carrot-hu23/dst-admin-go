package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initClusterRouter(router *gin.RouterGroup) {

	clusterApi := api.ClusterApi{}
	zoneApi := api.ZoneApi{}

	cluster := router.Group("/api/cluster")
	{
		cluster.GET("", clusterApi.GetClusterList)
		cluster.GET("/detail/:id", clusterApi.GetCluster)
		cluster.GET("/restart", clusterApi.RestartCluster)
		cluster.POST("", clusterApi.CreateCluster)
		cluster.PUT("", clusterApi.UpdateCluster)
		cluster.DELETE("", clusterApi.DeleteCluster)

		cluster.PUT("/container", clusterApi.UpdateClusterContainer)

		cluster.GET("/kami", clusterApi.GetKamiList)
		cluster.GET("/kami/export", clusterApi.ExportKamiList)

		cluster.GET("/zone", zoneApi.GetZone)
		cluster.POST("/zone", zoneApi.CreateZone)
		cluster.PUT("/zone", zoneApi.UpdateZone)
		cluster.DELETE("/zone", zoneApi.DeleteZone)

	}

	activate := router.Group("/activate")
	{
		activate.GET("/:id", clusterApi.GetCluster)
		activate.POST("/bind", clusterApi.BindCluster)
	}

}

package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initClusterRouter(router *gin.RouterGroup) {

	clusterApi := api.ClusterApi{}
	zoneApi := api.ZoneApi{}
	queueApi := api.QueueApi{}
	levelTemplateApi := api.LevelTemplateApi{}

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

		cluster.GET("/queue", queueApi.GetQueue)
		cluster.POST("/queue", queueApi.CreateQueue)
		cluster.PUT("/queue", queueApi.UpdateQueue)
		cluster.DELETE("/queue", queueApi.DeleteQueue)

		cluster.POST("/zone/queue/bind", queueApi.BindQueue2Zone)
		cluster.POST("/zone/queue/unbind", queueApi.UnbindQueueFromZone)
		cluster.GET("/zone/queue", queueApi.GetQueuesByZone)

		cluster.GET("/level/template", levelTemplateApi.GetLevelTemplate)
		cluster.POST("/level/template", levelTemplateApi.CreateLevelTemplate)
		cluster.PUT("/level/template", levelTemplateApi.UpdateTemplate)
		cluster.DELETE("/level/template", levelTemplateApi.DeleteTemplate)

	}

	activate := router.Group("/activate")
	{
		activate.GET("/:id", clusterApi.GetCluster)
		activate.POST("/bind", clusterApi.BindCluster)
	}

}

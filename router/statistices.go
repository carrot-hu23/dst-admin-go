package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initStatisticsRouter(router *gin.RouterGroup) {

	statisticsApi := api.StatisticsApi{}
	statistics := router.Group("/api/statistics")
	{
		statistics.GET("/active/user", statisticsApi.CountActiveUser)
		statistics.GET("/top/death", statisticsApi.TopDeaths)
		statistics.GET("/top/login", statisticsApi.TopUserLoginimes)
		statistics.GET("/top/active", statisticsApi.TopUserActiveTimes)

		statistics.GET("/rate/role", statisticsApi.CountRoleRate)
	}

}

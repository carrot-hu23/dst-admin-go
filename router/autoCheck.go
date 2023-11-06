package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initAutoCheck(router *gin.RouterGroup) {

	autoCheckApi := api.AutoCheckApi{}
	autoCheck := router.Group("/api/auto/check")
	{
		autoCheck.GET("/status", autoCheckApi.GetAutoCheckStatus)
		autoCheck.GET("/master", autoCheckApi.EnableAutoCheckMasterRun)
		autoCheck.GET("/caves", autoCheckApi.EnableAutoCheckCavesRun)
		autoCheck.GET("/version", autoCheckApi.EnableAutoCheckUpdateVersion)
		autoCheck.GET("/master/mod", autoCheckApi.EnableAutoCheckMasterMod)
		autoCheck.GET("/caves/mod", autoCheckApi.EnableAutoCheckCavesMod)

		autoCheck.GET("", autoCheckApi.GetAutoCheck)
		autoCheck.POST("", autoCheckApi.SaveAutoCheck)
	}

	autoCheck2 := router.Group("/api/auto/check2")
	{
		autoCheck2.GET("", autoCheckApi.GetAutoCheckList2)
		autoCheck2.POST("", autoCheckApi.SaveAutoCheck2)
	}

}

package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initAutoCheck(router *gin.RouterGroup) {

	autoCheckApi := api.AutoCheckApi{}

	autoCheck2 := router.Group("/api/auto/check2")
	{
		autoCheck2.GET("", autoCheckApi.GetAutoCheckList2)
		autoCheck2.POST("", autoCheckApi.SaveAutoCheck2)
	}

}

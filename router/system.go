package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initSystemRouter(router *gin.RouterGroup) {

	systemApi := api.SystemApi{}
	group := router.Group("/api/system")
	{
		group.GET("/setting", systemApi.GetSystemSetting)
		group.POST("/setting", systemApi.SaveSystemSetting)
		group.GET("/setting/install/steamcmd", systemApi.InstallSteamCmd)
	}

}

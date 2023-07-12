package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initFirstRouter(router *gin.RouterGroup) {

	fistApi := api.InitApi{}
	router.GET("/api/init", fistApi.CheckIsFirst)
	router.POST("/api/init", fistApi.InitFirst)
	// router.GET("/api/install/steamcmd", fistApi.InstallSteamCmd)

	installSteamCmdApi := api.InstallSteamCmd{}
	router.GET("/api/install/steamcmd", installSteamCmdApi.InstallSteamCmd)

}

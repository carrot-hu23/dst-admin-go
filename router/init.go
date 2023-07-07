package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initFirstRouter(router *gin.RouterGroup) {

	fistApi := api.InitApi{}
	router.GET("/api/initConfig", fistApi.CheckIsFirst)
	router.POST("/api/initConfig", fistApi.InitFirst)
	// router.GET("/api/install/steamcmd", fistApi.InstallSteamCmd)

	installSteamCmdApi := api.InstallSteamCmd{}
	router.GET("/api/install/steamcmd", installSteamCmdApi.InstallSteamCmd)

}

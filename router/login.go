package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initLoginRouter(router *gin.RouterGroup) {

	loginApi := api.LoginApi{}
	router.POST("/api/login", loginApi.Login)
	router.GET("/api/logout", loginApi.Logout)
	router.POST("/api/change/password", loginApi.ChangePassword)
	router.GET("/api/user", loginApi.GetUserInfo)
	router.POST("/api/user", loginApi.UpdateUserInfo)
}

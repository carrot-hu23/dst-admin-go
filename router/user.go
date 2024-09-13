package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initUserRouter(router *gin.RouterGroup) {

	userApi := api.UserApi{}

	user := router.Group("/api/user/account")
	{
		user.GET("", userApi.QueryUserList)
		user.POST("", userApi.CreateUser)
		user.PUT("", userApi.UpdateUser)
		user.DELETE("", userApi.DeleteUser)
	}

	userCluster := router.Group("/api/user/account/cluster")
	{
		userCluster.GET("", userApi.GetUserClusterList)
		userCluster.POST("", userApi.AddUserCluster)
		userCluster.DELETE("", userApi.RemoveUserCluster)
		userCluster.PUT("", userApi.UpdateUserAllow)

		userCluster.GET("/permission", userApi.GetUserCluster)
	}
}

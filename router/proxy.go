package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initProxyRouter(router *gin.RouterGroup) {

	proxyApi := api.ProxyApi{}
	router.Any("/app/:name/*path", proxyApi.NewProxy)

	proxyApp := router.Group("/api/proxy")
	{
		proxyApp.GET("", proxyApi.GetProxyEntity)
		proxyApp.POST("", proxyApi.CreateProxyEntity)
		proxyApp.PUT("", proxyApi.UpdateProxyEntity)
		proxyApp.DELETE("", proxyApi.DeleteProxyEntity)
	}

}

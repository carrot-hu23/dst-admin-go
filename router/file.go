package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initFile(router *gin.RouterGroup) {

	fileApi := api.FileApi{}
	file := router.Group("/api/file")
	{
		file.POST("/ugc/upload", fileApi.UploadUgcMods)
	}

}

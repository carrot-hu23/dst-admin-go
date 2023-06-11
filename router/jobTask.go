package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initJobTaskRouter(router *gin.RouterGroup) {

	taskApi := api.JobTaskApi{}
	task := router.Group("/api/task")
	{
		task.GET("", taskApi.GetJobTaskList)
		task.POST("", taskApi.AddJobTask)
		task.DELETE("", taskApi.DeleteJobTask)
	}

}

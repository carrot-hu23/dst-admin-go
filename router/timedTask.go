package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initTimedTaskRouter(router *gin.RouterGroup) {

	taskApi := api.TimedTaskApi{}
	task := router.Group("/api/task")
	{
		task.GET("", taskApi.GetJobTaskList)
		task.POST("", taskApi.AddJobTask)
		task.DELETE("", taskApi.DeleteJobTask)
	}

}

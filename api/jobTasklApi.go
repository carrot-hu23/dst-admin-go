package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/schedule"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type JobTaskApi struct {
}

func (j *JobTaskApi) GetInstructList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: schedule.ScheduleSingleton.GetInstructList(),
	})
}

func (j *JobTaskApi) GetJobTaskList(ctx *gin.Context) {

	jobs := schedule.ScheduleSingleton.GetJobs()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: jobs,
	})
}

func (j *JobTaskApi) AddJobTask(ctx *gin.Context) {

	jobTask := &model.JobTask{}
	if err := ctx.ShouldBindJSON(jobTask); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	db := database.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, vo.Response{
				Code: 500,
				Msg:  "create task err",
				Data: nil,
			})
		}
	}()

	tx.Create(jobTask)
	task := schedule.Task{
		Id:          jobTask.ID,
		Corn:        jobTask.Cron,
		F:           schedule.StrategyMap[jobTask.Category].Execute,
		ClusterName: jobTask.ClusterName,
	}
	schedule.ScheduleSingleton.AddJob(task)
	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (j *JobTaskApi) DeleteJobTask(ctx *gin.Context) {

	jobId, _ := strconv.Atoi(ctx.DefaultQuery("jobId", "0"))
	schedule.ScheduleSingleton.DeleteJob(jobId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/schedule"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type TimedTaskApi struct {
}

func (j *TimedTaskApi) GetInstructList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: schedule.ScheduleSingleton.GetInstructList(),
	})
}

func (j *TimedTaskApi) GetJobTaskList(ctx *gin.Context) {

	jobs := schedule.ScheduleSingleton.GetJobs()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: jobs,
	})
}

func (j *TimedTaskApi) AddJobTask(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	jobTask := &model.JobTask{}
	if err := ctx.ShouldBindJSON(jobTask); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	if jobTask.ClusterName == "" {
		jobTask.ClusterName = cluster.ClusterName
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
		Id:           jobTask.ID,
		Corn:         jobTask.Cron,
		F:            schedule.StrategyMap[jobTask.Category].Execute,
		ClusterName:  jobTask.ClusterName,
		Announcement: jobTask.Announcement,
		Sleep:        jobTask.Sleep,
		Times:        jobTask.Times,
	}
	schedule.ScheduleSingleton.AddJob(task)
	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (j *TimedTaskApi) DeleteJobTask(ctx *gin.Context) {

	jobId, _ := strconv.Atoi(ctx.DefaultQuery("jobId", "0"))
	log.Println("jobid: ", jobId)
	schedule.ScheduleSingleton.DeleteJob(jobId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

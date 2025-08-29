package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/schedule"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

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

	// 绑定 JSON
	if err := ctx.ShouldBindJSON(jobTask); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Response{
			Code: 400,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 填充 clusterName
	if jobTask.ClusterName == "" {
		jobTask.ClusterName = cluster.ClusterName
	}

	// 校验 cron
	if jobTask.Cron == "" {
		ctx.JSON(http.StatusBadRequest, vo.Response{
			Code: 400,
			Msg:  "cron 表达式不能为空",
		})
		return
	}

	// 使用 robfig/cron 校验格式
	if _, err := cron.ParseStandard(jobTask.Cron); err != nil {
		ctx.JSON(http.StatusBadRequest, vo.Response{
			Code: 400,
			Msg:  "cron 表达式格式错误: " + err.Error(),
		})
		return
	}

	db := database.DB
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, vo.Response{
				Code: 500,
				Msg:  "创建任务异常",
			})
		}
	}()

	// 创建任务记录
	if err := tx.Create(jobTask).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, vo.Response{
			Code: 500,
			Msg:  "创建任务失败: " + err.Error(),
		})
		return
	}

	// 调度任务
	task := schedule.Task{
		Id:           jobTask.ID,
		Corn:         jobTask.Cron,
		F:            schedule.StrategyMap[jobTask.Category].Execute,
		ClusterName:  jobTask.ClusterName,
		LevelName:    jobTask.Uuid,
		Announcement: jobTask.Announcement,
		Sleep:        jobTask.Sleep,
		Times:        jobTask.Times,
	}
	schedule.ScheduleSingleton.AddJob(task)

	tx.Commit()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
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

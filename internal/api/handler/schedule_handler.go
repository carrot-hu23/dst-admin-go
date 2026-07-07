package handler

import (
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/schedule"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type ScheduleHandler struct {
	scheduleService *schedule.ScheduleService
}

func NewScheduleHandler(scheduleService *schedule.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

func (h *ScheduleHandler) RegisterRoute(router *gin.RouterGroup) {
	task := router.Group("/api/task")
	{
		task.GET("", h.GetJobTaskList)
		task.POST("", h.AddJobTask)
		task.DELETE("", h.DeleteJobTask)
		task.GET("/instruct", h.GetInstructList)
	}
}

// GetInstructList 获取指令列表
// @Summary 获取指令列表
// @Description 获取可用的定时任务指令列表
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/task/instruct [get]
func (h *ScheduleHandler) GetInstructList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: h.scheduleService.GetInstructList(),
	})
}

// GetJobTaskList 获取定时任务列表
// @Summary 获取定时任务列表
// @Description 获取所有定时任务的列表
// @Tags task
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/task [get]
func (h *ScheduleHandler) GetJobTaskList(ctx *gin.Context) {
	jobs := h.scheduleService.GetJobs()

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: jobs,
	})
}

// AddJobTask 添加定时任务
// @Summary 添加定时任务
// @Description 创建新的定时任务
// @Tags task
// @Accept json
// @Produce json
// @Param data body model.JobTask true "任务信息"
// @Success 200 {object} response.Response
// @Router /api/task [post]
func (h *ScheduleHandler) AddJobTask(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	jobTask := &model.JobTask{}

	// 绑定 JSON
	if err := ctx.ShouldBindJSON(jobTask); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 400,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	// 填充 clusterName
	if jobTask.ClusterName == "" {
		jobTask.ClusterName = clusterName
	}

	// 校验 cron
	if jobTask.Cron == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "cron 表达式不能为空",
		})
		return
	}

	// 使用 robfig/cron 校验格式
	if _, err := cron.ParseStandard(jobTask.Cron); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "cron 表达式格式错误: " + err.Error(),
		})
		return
	}

	// 创建任务
	if err := h.scheduleService.CreateJobTask(jobTask); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "创建任务失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
	})
}

// DeleteJobTask 删除定时任务
// @Summary 删除定时任务
// @Description 根据任务ID删除定时任务
// @Tags task
// @Accept json
// @Produce json
// @Param jobId query int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/task [delete]
func (h *ScheduleHandler) DeleteJobTask(ctx *gin.Context) {
	jobID, _ := strconv.Atoi(ctx.DefaultQuery("jobId", "0"))
	log.Println("jobid: ", jobID)

	if err := h.scheduleService.DeleteJob(jobID); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "删除任务失败: " + err.Error(),
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

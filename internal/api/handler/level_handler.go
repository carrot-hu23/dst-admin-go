package handler

import (
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/level"
	"dst-admin-go/internal/service/levelConfig"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LevelHandler struct {
	levelService *level.LevelService
}

func NewLevelHandler(levelService *level.LevelService) *LevelHandler {
	return &LevelHandler{
		levelService: levelService,
	}
}

func (h *LevelHandler) RegisterRoute(router *gin.RouterGroup) {
	level := router.Group("/api/cluster/level")
	{
		level.GET("", h.GetLevelList)
		level.POST("", h.CreateLevel)
		level.DELETE("", h.DeleteLevel)
		level.PUT("", h.UpdateLevels)
	}
}

// GetLevelList 获取世界列表
// @Summary 获取世界列表
// @Description 获取当前集群的所有世界(等级)列表
// @Tags level
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]levelConfig.LevelInfo}
// @Router /api/cluster/level [get]
func (h *LevelHandler) GetLevelList(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	levels := h.levelService.GetLevelList(clusterName)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "get level list success",
		Data: levels,
	})
}

// UpdateLevel 更新世界
// @Summary 更新世界
// @Description 更新指定世界的配置信息
// @Tags level
// @Accept json
// @Produce json
// @Param level body levelConfig.LevelInfo true "世界配置信息"
// @Success 200 {object} response.Response
// @Router /api/cluster/level [put]
func (h *LevelHandler) UpdateLevel(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)

	var world levelConfig.LevelInfo
	if err := ctx.BindJSON(&world); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 400,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	err := h.levelService.UpdateLevel(clusterName, &world)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "update level failed",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "update level success",
		Data: nil,
	})
}

// CreateLevel 创建世界
// @Summary 创建世界
// @Description 创建一个新的世界(等级)
// @Tags level
// @Accept json
// @Produce json
// @Param level body levelConfig.LevelInfo true "世界配置信息"
// @Success 200 {object} response.Response
// @Router /api/cluster/level [post]
func (h *LevelHandler) CreateLevel(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)

	var world levelConfig.LevelInfo
	if err := ctx.BindJSON(&world); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 400,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	err := h.levelService.CreateLevel(clusterName, &world)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "create level failed",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "create level success",
		Data: world,
	})
}

// DeleteLevel 删除世界
// @Summary 删除世界
// @Description 删除指定的世界
// @Tags level
// @Accept json
// @Produce json
// @Param levelName query string true "世界名称"
// @Success 200 {object} response.Response
// @Router /api/cluster/level [delete]
func (h *LevelHandler) DeleteLevel(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	levelName := ctx.Query("levelName")

	err := h.levelService.DeleteLevel(clusterName, levelName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "delete level failed",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "delete level success",
		Data: nil,
	})
}

// UpdateLevels 批量更新世界
// @Summary 批量更新世界
// @Description 批量更新多个世界的配置信息
// @Tags level
// @Accept json
// @Produce json
// @Param levels body []levelConfig.LevelInfo true "世界配置信息列表"
// @Success 200 {object} response.Response
// @Router /api/cluster/level [put]
func (h *LevelHandler) UpdateLevels(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	var payload struct {
		Levels []levelConfig.LevelInfo `json:"levels"`
	}
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 400,
			Msg:  "invalid request",
			Data: nil,
		})
		return
	}

	err := h.levelService.UpdateLevels(clusterName, payload.Levels)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "update levels failed",
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "update levels success",
		Data: nil,
	})
}

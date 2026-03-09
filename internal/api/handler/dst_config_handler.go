package handler

import (
	"dst-admin-go/internal/collect"
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/dstConfig"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DstConfigHandler struct {
	dstConfig dstConfig.Config
	archive   *archive.PathResolver
}

func NewDstConfigHandler(dstConfig dstConfig.Config, archive *archive.PathResolver) *DstConfigHandler {
	return &DstConfigHandler{
		dstConfig: dstConfig,
		archive:   archive,
	}
}

func (h *DstConfigHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("api/dst/config", h.GetDstConfig)
	router.POST("api/dst/config", h.SaveDstConfig)
}

// GetDstConfig 获取房间 dst 配置 swagger 注释
// @Summary 获取房间 dst 配置
// @Description 获取房间 dst 配置
// @Tags dstConfig
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=dstConfig.DstConfig}
// @Router /api/game/dst/config [get]
func (h *DstConfigHandler) GetDstConfig(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)

	config, err := h.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "failed to get dst config: " + err.Error(),
			Data: nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: config,
	})
}

// SaveDstConfig 保存房间 dst 配置 swagger 注释
// @Summary 保存房间 dst 配置
// @Description 保存房间 dst 配置
// @Tags dstConfig
// @Accept json
// @Produce json
// @Param config body dstConfig.DstConfig true "dst 配置"
// @Success 200 {object} response.Response
// @Router /api/game/dst/config [post]
func (h *DstConfigHandler) SaveDstConfig(ctx *gin.Context) {
	config := dstConfig.DstConfig{}
	if err := ctx.ShouldBindJSON(&config); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	err := h.dstConfig.SaveDstConfig(context.GetClusterName(ctx), config)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code: 500,
			Msg:  "failed to save dst config: " + err.Error(),
			Data: nil,
		})
		return
	}
	clusterPath := h.archive.ClusterPath(config.Cluster)
	collect.Collector.ReCollect(clusterPath, config.Cluster)
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "DstConfig saved successfully",
		Data: nil,
	})
}

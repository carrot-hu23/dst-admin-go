package handler

import (
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/mod"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ModHandler struct {
	modService *mod.ModService
	dstConfig  dstConfig.Config
}

func NewModHandler(modService *mod.ModService, dstConfig dstConfig.Config) *ModHandler {
	return &ModHandler{
		modService: modService,
		dstConfig:  dstConfig,
	}
}

func (h *ModHandler) RegisterRoute(router *gin.RouterGroup) {
	modGroup := router.Group("/api/mod")
	{
		modGroup.GET("/search", h.SearchModList)
		modGroup.GET("/:modId", h.GetModInfo)
		modGroup.PUT("/:modId", h.UpdateMod)
		modGroup.GET("", h.GetMyModList)
		modGroup.DELETE("/:modId", h.DeleteMod)
		modGroup.DELETE("/setup/workshop", h.DeleteSetupWorkshop)
		modGroup.GET("/modinfo/:modId", h.GetModInfoFile)
		modGroup.POST("/modinfo", h.SaveModInfoFile)
		modGroup.POST("/modinfo/file", h.AddModInfoFile)
		modGroup.PUT("/modinfo", h.UpdateAllModInfos)
		modGroup.GET("/ugc/acf", h.GetUgcModAcf)
		modGroup.DELETE("/ugc", h.DeleteUgcModFile)
	}
}

// SearchModList 搜索mod列表
// @Summary 搜索mod列表
// @Description 搜索Steam创意工坊中的模组
// @Tags mod
// @Accept json
// @Produce json
// @Param text query string false "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param lang query string false "语言" default(zh)
// @Success 200 {object} response.Response{}
// @Router /api/mod/search [get]
func (h *ModHandler) SearchModList(ctx *gin.Context) {
	text := ctx.Query("text")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	lang := ctx.DefaultQuery("lang", "zh")

	data, err := h.modService.SearchModList(text, page, size, lang)
	if err != nil {
		response.FailWithMessage("搜索mod失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(data, ctx)
}

// GetModInfo 获取mod信息
// @Summary 获取mod信息
// @Description 根据modId获取模组详细信息
// @Tags mod
// @Accept json
// @Produce json
// @Param modId path string true "模组ID"
// @Param lang query string false "语言" default(zh)
// @Success 200 {object} response.Response
// @Router /api/mod/{modId} [get]
func (h *ModHandler) GetModInfo(ctx *gin.Context) {
	modId := ctx.Param("modId")
	lang := ctx.DefaultQuery("lang", "zh")
	clusterName := context.GetClusterName(ctx)
	modinfo, err := h.modService.SubscribeModByModId(clusterName, modId, lang)
	if err != nil {
		response.FailWithMessage("模组下载失败: "+err.Error(), ctx)
		return
	}

	var modConfig map[string]interface{}
	_ = json.Unmarshal([]byte(modinfo.ModConfig), &modConfig)

	modData := map[string]interface{}{
		"auth":          modinfo.Auth,
		"consumer_id":   modinfo.ConsumerAppid,
		"creator_appid": modinfo.CreatorAppid,
		"description":   modinfo.Description,
		"file_url":      modinfo.FileUrl,
		"modid":         modinfo.Modid,
		"img":           modinfo.Img,
		"last_time":     modinfo.LastTime,
		"name":          modinfo.Name,
		"v":             modinfo.V,
		"mod_config":    modConfig,
		"update":        modinfo.Update,
	}

	response.OkWithData(modData, ctx)
}

// GetMyModList 获取我的mod列表
// @Summary 获取我的mod列表
// @Description 获取已订阅的模组列表
// @Tags mod
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/mod [get]
func (h *ModHandler) GetMyModList(ctx *gin.Context) {
	modInfos, err := h.modService.GetMyModList()
	if err != nil {
		response.FailWithMessage("获取模组列表失败: "+err.Error(), ctx)
		return
	}

	var modDataList []map[string]interface{}
	for _, modinfo := range modInfos {
		var modConfig map[string]interface{}
		_ = json.Unmarshal([]byte(modinfo.ModConfig), &modConfig)

		modData := map[string]interface{}{
			"auth":          modinfo.Auth,
			"consumer_id":   modinfo.ConsumerAppid,
			"creator_appid": modinfo.CreatorAppid,
			"description":   modinfo.Description,
			"file_url":      modinfo.FileUrl,
			"modid":         modinfo.Modid,
			"img":           modinfo.Img,
			"last_time":     modinfo.LastTime,
			"name":          modinfo.Name,
			"v":             modinfo.V,
			"mod_config":    modConfig,
			"update":        modinfo.Update,
		}
		modDataList = append(modDataList, modData)
	}

	response.OkWithData(modDataList, ctx)
}

// UpdateAllModInfos 批量更新模组信息
// @Summary 批量更新模组信息
// @Description 批量更新所有已订阅模组的信息
// @Tags mod
// @Accept json
// @Produce json
// @Param lang query string false "语言" default(zh)
// @Success 200 {object} response.Response
// @Router /api/mod/modinfo [put]
func (h *ModHandler) UpdateAllModInfos(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)
	lang := ctx.DefaultQuery("lang", "zh")

	err := h.modService.UpdateAllModInfos(clusterName, lang)
	if err != nil {
		response.FailWithMessage("更新失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("更新成功", ctx)
}

// DeleteMod 删除模组
// @Summary 删除模组
// @Description 根据modId删除模组
// @Tags mod
// @Accept json
// @Produce json
// @Param modId path string true "模组ID"
// @Success 200 {object} response.Response
// @Router /api/mod/{modId} [delete]
func (h *ModHandler) DeleteMod(ctx *gin.Context) {
	modId := ctx.Param("modId")
	clusterName := context.GetClusterName(ctx)

	err := h.modService.DeleteMod(clusterName, modId)
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(modId, ctx)
}

// DeleteSetupWorkshop 删除workshop文件
// @Summary 删除workshop文件
// @Description 删除所有workshop模组文件
// @Tags mod
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/mod/setup/workshop [delete]
func (h *ModHandler) DeleteSetupWorkshop(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)

	err := h.modService.DeleteSetupWorkshop(clusterName)
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

// GetModInfoFile 获取模组配置文件
// @Summary 获取模组配置文件
// @Description 根据modId获取模组配置文件内容
// @Tags mod
// @Accept json
// @Produce json
// @Param modId path string true "模组ID"
// @Success 200 {object} response.Response{data=mod.ModInfo}
// @Router /api/mod/modinfo/{modId} [get]
func (h *ModHandler) GetModInfoFile(ctx *gin.Context) {
	modId := ctx.Param("modId")

	modInfo, err := h.modService.GetModByModId(modId)
	if err != nil {
		response.FailWithMessage("获取模组信息失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(modInfo, ctx)
}

// SaveModInfoFile 保存模组配置文件
// @Summary 保存模组配置文件
// @Description 保存模组配置信息
// @Tags mod
// @Accept json
// @Produce json
// @Param data body mod.ModInfo true "模组信息"
// @Success 200 {object} response.Response{data=mod.ModInfo}
// @Router /api/mod/modinfo [post]
func (h *ModHandler) SaveModInfoFile(ctx *gin.Context) {
	var modInfo model.ModInfo
	err := ctx.ShouldBindJSON(&modInfo)
	if err != nil {
		response.FailWithMessage("参数解析失败: "+err.Error(), ctx)
		return
	}

	err = h.modService.SaveModInfo(&modInfo)
	if err != nil {
		response.FailWithMessage("保存失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(modInfo, ctx)
}

// UpdateMod 更新模组
// @Summary 更新模组
// @Description 根据modId更新模组
// @Tags mod
// @Accept json
// @Produce json
// @Param modId path string true "模组ID"
// @Param lang query string false "语言" default(zh)
// @Success 200 {object} response.Response
// @Router /api/mod/{modId} [put]
func (h *ModHandler) UpdateMod(ctx *gin.Context) {
	modId := ctx.Param("modId")
	clusterName := context.GetClusterName(ctx)
	lang := ctx.DefaultQuery("lang", "zh")

	// 删除旧数据
	err := h.modService.DeleteMod(clusterName, modId)
	if err != nil {
		log.Println("删除旧模组失败:", err)
	}

	// 重新下载
	modinfo, err := h.modService.SubscribeModByModId(clusterName, modId, lang)
	if err != nil {
		response.FailWithMessage("模组更新失败: "+err.Error(), ctx)
		return
	}

	var modConfig map[string]interface{}
	_ = json.Unmarshal([]byte(modinfo.ModConfig), &modConfig)

	modData := map[string]interface{}{
		"auth":          modinfo.Auth,
		"consumer_id":   modinfo.ConsumerAppid,
		"creator_appid": modinfo.CreatorAppid,
		"description":   modinfo.Description,
		"file_url":      modinfo.FileUrl,
		"modid":         modinfo.Modid,
		"img":           modinfo.Img,
		"last_time":     modinfo.LastTime,
		"name":          modinfo.Name,
		"v":             modinfo.V,
		"mod_config":    modConfig,
		"update":        modinfo.Update,
	}

	response.OkWithData(modData, ctx)
}

// AddModInfoFile 手动添加模组
// @Summary 手动添加模组
// @Description 手动添加模组配置文件
// @Tags mod
// @Accept json
// @Produce json
// @Param data body object true "模组信息" example({"workshopId":"123456","modinfo":"模组配置内容"})
// @Param lang query string false "语言" default(zh)
// @Success 200 {object} response.Response
// @Router /api/mod/modinfo/file [post]
func (h *ModHandler) AddModInfoFile(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	lang := ctx.DefaultQuery("lang", "zh")

	var payload struct {
		WorkshopId string `json:"workshopId"`
		Modinfo    string `json:"modinfo"`
	}

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		response.FailWithMessage("参数解析失败: "+err.Error(), ctx)
		return
	}

	if payload.WorkshopId == "" {
		response.FailWithMessage("workshopId不能为空", ctx)
		return
	}

	// 获取配置
	config, err := h.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		response.FailWithMessage("获取配置失败: "+err.Error(), ctx)
		return
	}

	err = h.modService.AddModInfo(clusterName, lang, payload.WorkshopId, payload.Modinfo, config.Mod_download_path)
	if err != nil {
		response.FailWithMessage("添加模组失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("添加成功", ctx)
}

// GetUgcModAcf 获取UGC mod acf文件信息
// @Summary 获取UGC mod acf文件信息
// @Description 获取UGC模组的ACF文件信息
// @Tags mod
// @Accept json
// @Produce json
// @Param levelName query string true "世界名称"
// @Success 200 {object} response.Response
// @Router /api/mod/ugc/acf [get]
func (h *ModHandler) GetUgcModAcf(ctx *gin.Context) {

	levelName := ctx.Query("levelName")
	clusterName := context.GetClusterName(ctx)

	workshopItemDetails, err := h.modService.GetUgcModInfo(clusterName, levelName)
	if err != nil {
		response.FailWithMessage("获取UGC模组信息失败: "+err.Error(), ctx)
		return
	}

	response.OkWithData(workshopItemDetails, ctx)
}

// DeleteUgcModFile 删除UGC模组文件
// @Summary 删除UGC模组文件
// @Description 删除UGC模组文件
// @Tags mod
// @Accept json
// @Produce json
// @Param levelName query string true "世界名称"
// @Param workshopId query string true "WorkshopID"
// @Success 200 {object} response.Response
// @Router /api/mod/ugc [delete]
func (h *ModHandler) DeleteUgcModFile(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	levelName := ctx.Query("levelName")
	workshopId := ctx.Query("workshopId")

	err := h.modService.DeleteUgcModFile(clusterName, levelName, workshopId)
	if err != nil {
		response.FailWithMessage("删除失败: "+err.Error(), ctx)
		return
	}

	response.OkWithMessage("删除成功", ctx)
}

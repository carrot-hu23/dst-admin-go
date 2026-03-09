package handler

import (
	"dst-admin-go/internal/database"
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/backup"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	backupService *backup.BackupService
}

func NewBackupHandler(backupService *backup.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

func (h *BackupHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/game/backup", h.GetBackupList)
	router.POST("/api/game/backup", h.CreateBackup)
	router.DELETE("/api/game/backup", h.DeleteBackup)
	router.PUT("/api/game/backup", h.RenameBackup)
	router.GET("/api/game/backup/download", h.DownloadBackup)
	router.POST("/api/game/backup/upload", h.UploadBackup)
	router.GET("/backup/restore", h.RestoreBackup)
	router.POST("/api/game/backup/snapshot/setting", h.SaveBackupSnapshotsSetting)
	router.GET("/api/game/backup/snapshot/setting", h.GetBackupSnapshotsSetting)
	router.GET("/api/game/backup/snapshot/list", h.BackupSnapshotsList)
}

// DeleteBackup 删除备份
// @Summary 删除备份
// @Description 删除指定的备份文件
// @Tags backup
// @Accept json
// @Produce json
// @Param fileNames body []string true "要删除的文件名列表"
// @Success 200 {object} response.Response
// @Router /api/game/backup [delete]
func (h *BackupHandler) DeleteBackup(ctx *gin.Context) {
	var body struct {
		FileNames []string `json:"fileNames"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	h.backupService.DeleteBackup(ctx, body.FileNames)
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "delete backups success",
		Data: nil,
	})
}

// DownloadBackup 下载备份
// @Summary 下载备份
// @Description 下载指定的备份文件
// @Tags backup
// @Accept json
// @Produce json
// @Param backupName query string true "备份文件名"
// @Success 200 {object} response.Response
// @Router /api/game/backup/download [get]
func (h *BackupHandler) DownloadBackup(ctx *gin.Context) {
	h.backupService.DownloadBackup(ctx)
}

// GetBackupList 获取备份列表
// @Summary 获取备份列表
// @Description 获取当前集群的所有备份文件列表
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/backup [get]
func (h *BackupHandler) GetBackupList(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "get backup list success",
		Data: h.backupService.GetBackupList(clusterName),
	})
}

// RenameBackup 重命名备份
// @Summary 重命名备份
// @Description 重命名指定的备份文件
// @Tags backup
// @Accept json
// @Produce json
// @Param request body object true "请求体" schema-example({"fileName": "old_name", "newName": "new_name"})
// @Success 200 {object} response.Response
// @Router /api/game/backup [put]
func (h *BackupHandler) RenameBackup(ctx *gin.Context) {

	var body struct {
		FileName string `json:"fileName"`
		NewName  string `json:"newName"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	h.backupService.RenameBackup(ctx, body.FileName, body.NewName)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "rename backup success",
		Data: nil,
	})
}

// RestoreBackup 恢复备份
// @Summary 恢复备份
// @Description 从备份文件恢复游戏存档
// @Tags backup
// @Accept json
// @Produce json
// @Param backupName query string true "备份文件名"
// @Success 200 {object} response.Response
// @Router /backup/restore [get]
func (h *BackupHandler) RestoreBackup(ctx *gin.Context) {
	backupName := ctx.Query("backupName")

	h.backupService.RestoreBackup(ctx, backupName)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "restore backup success",
		Data: nil,
	})
}

// UploadBackup 上传备份
// @Summary 上传备份
// @Description 上传新的备份文件
// @Tags backup
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/backup/upload [post]
func (h *BackupHandler) UploadBackup(ctx *gin.Context) {

	h.backupService.UploadBackup(ctx)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "upload backup success",
		Data: nil,
	})
}

// CreateBackup 创建备份
// @Summary 创建备份
// @Description 创建新的游戏存档备份
// @Tags backup
// @Accept json
// @Produce json
// @Param backupName query string false "备份名称"
// @Success 200 {object} response.Response
// @Router /api/game/backup [post]
func (h *BackupHandler) CreateBackup(ctx *gin.Context) {
	var body struct {
		BackupName string `json:"backupName"`
	}
	if err := ctx.ShouldBind(&body); err != nil {
		body.BackupName = ""
	}
	clusterName := context.GetClusterName(ctx)
	h.backupService.CreateBackup(clusterName, body.BackupName)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "create backup success",
		Data: nil,
	})
}

// SaveBackupSnapshotsSetting 保存快照设置
// @Summary 保存快照设置
// @Description 保存自动快照备份的设置
// @Tags backup
// @Accept json
// @Produce json
// @Param setting body model.BackupSnapshot true "快照设置"
// @Success 200 {object} response.Response
// @Router /api/game/backup/snapshot/setting [post]
func (h *BackupHandler) SaveBackupSnapshotsSetting(ctx *gin.Context) {

	var backupSnapshot model.BackupSnapshot
	var oldBackupSnapshot model.BackupSnapshot
	err := ctx.ShouldBind(&backupSnapshot)
	if err != nil {
		log.Panicln("参数错误", err)
	}
	db := database.Db
	db.First(&oldBackupSnapshot)

	oldBackupSnapshot.Enable = backupSnapshot.Enable
	oldBackupSnapshot.Interval = backupSnapshot.Interval
	oldBackupSnapshot.MaxSnapshots = backupSnapshot.MaxSnapshots
	oldBackupSnapshot.IsCSave = backupSnapshot.IsCSave

	db.Save(&oldBackupSnapshot)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: oldBackupSnapshot,
	})
}

// GetBackupSnapshotsSetting 获取快照设置
// @Summary 获取快照设置
// @Description 获取自动快照备份的设置
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/backup/snapshot/setting [get]
func (h *BackupHandler) GetBackupSnapshotsSetting(ctx *gin.Context) {

	var oldBackupSnapshot model.BackupSnapshot
	db := database.Db
	db.First(&oldBackupSnapshot)

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: oldBackupSnapshot,
	})
}

// BackupSnapshotsList 获取快照列表
// @Summary 获取快照列表
// @Description 获取当前集群的所有快照备份列表
// @Tags backup
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/game/backup/snapshot/list [get]
func (h *BackupHandler) BackupSnapshotsList(ctx *gin.Context) {
	clusterName := context.GetClusterName(ctx)

	var snapshotBackupList []backup.BackupInfo
	backupList := h.backupService.GetBackupList(clusterName)
	for i := range backupList {
		name := backupList[i].FileName
		if strings.HasPrefix(name, "(snapshot)") && strings.Contains(name, clusterName) {
			snapshotBackupList = append(snapshotBackupList, backupList[i])
		}
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: snapshotBackupList,
	})
}

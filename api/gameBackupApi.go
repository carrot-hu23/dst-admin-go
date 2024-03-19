package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type GameBackUpApi struct {
}

var backupService = service.BackupService{}

func (g *GameBackUpApi) DeleteBackup(ctx *gin.Context) {
	var body struct {
		FileNames []string `json:"fileNames"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	backupService.DeleteBackup(ctx, body.FileNames)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "delete backups success",
		Data: nil,
	})
}

func (g *GameBackUpApi) DownloadBackup(ctx *gin.Context) {
	backupService.DownloadBackup(ctx)
}

func (g *GameBackUpApi) GetBackupList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "get backup list success",
		Data: backupService.GetBackupList(cluster.ClusterName),
	})
}

func (g *GameBackUpApi) RenameBackup(ctx *gin.Context) {

	var body struct {
		FileName string `json:"fileName"`
		NewName  string `json:"newName"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	backupService.RenameBackup(ctx, body.FileName, body.NewName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "rename backup success",
		Data: nil,
	})
}

func (g *GameBackUpApi) UploadBackup(ctx *gin.Context) {

	backupService.UploadBackup(ctx)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "upload backup success",
		Data: nil,
	})
}

func (g *GameBackUpApi) CreateBackup(ctx *gin.Context) {
	var body struct {
		BackupName string `json:"backupName"`
	}
	if err := ctx.ShouldBind(&body); err != nil {
		body.BackupName = ""
	}
	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupService.CreateBackup(cluster.ClusterName, body.BackupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "create backup success",
		Data: nil,
	})
}

func (g *GameBackUpApi) SaveBackupSnapshotsSetting(ctx *gin.Context) {

	var backupSnapshot model.BackupSnapshot
	var oldBackupSnapshot model.BackupSnapshot
	err := ctx.ShouldBind(&backupSnapshot)
	if err != nil {
		log.Panicln("参数错误", err)
	}
	db := database.DB
	db.First(&oldBackupSnapshot)

	oldBackupSnapshot.Enable = backupSnapshot.Enable
	oldBackupSnapshot.Interval = backupSnapshot.Interval
	oldBackupSnapshot.MaxSnapshots = backupSnapshot.MaxSnapshots
	oldBackupSnapshot.IsCSave = backupSnapshot.IsCSave

	db.Save(&oldBackupSnapshot)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: oldBackupSnapshot,
	})
}

func (g *GameBackUpApi) GetBackupSnapshotsSetting(ctx *gin.Context) {

	var oldBackupSnapshot model.BackupSnapshot
	db := database.DB
	db.First(&oldBackupSnapshot)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: oldBackupSnapshot,
	})
}

func (g *GameBackUpApi) BackupSnapshotsList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)

	var snapshotBackupList []vo.BackupVo
	backupList := backupService.GetBackupList(cluster.ClusterName)
	for i := range backupList {
		name := backupList[i].FileName
		if strings.HasPrefix(name, "(snapshot)") && strings.Contains(name, cluster.ClusterName) {
			snapshotBackupList = append(snapshotBackupList, backupList[i])
		}
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: snapshotBackupList,
	})
}

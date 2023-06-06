package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
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

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "get backup list success",
		Data: backupService.GetBackupList(ctx),
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

// UploadBackup TODO test
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
	backupService.CreateBackup(ctx, body.BackupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "create backup success",
		Data: nil,
	})
}

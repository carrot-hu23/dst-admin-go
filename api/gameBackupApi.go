package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

// /backup/deleteBackup
// DELETE /api/game/back/{backupName}
func DeleteBackup(ctx *gin.Context) {
	var body struct {
		FileNames []string `json:"fileNames"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	service.DeleteBackup(body.FileNames)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "delete backups success",
		Data: nil,
	})
}

// /backup/download
// GET /api/game/back/download/{backupName}
func DownloadBackup(ctx *gin.Context) {
	service.DownloadBackup(ctx)
}

// /backup/getBackupList
// GET /api/game/back/
func GetBackupList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "get backup list success",
		Data: service.GetBackupList(),
	})
}

// PUT /api/game/back/
func RenameBackup(ctx *gin.Context) {

	var body struct {
		FileName string `json:"fileName"`
		NewName  string `json:"NewName"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	service.RenameBackup(body.FileName, body.NewName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "rename backup success",
		Data: nil,
	})
}

// POST /api/game/back/
func UploadBackup(ctx *gin.Context) {

	service.UploadBackup(ctx)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "upload backup success",
		Data: nil,
	})
}

func CreateBackup(ctx *gin.Context) {
	var body struct {
		BackupName string `json:"backupName"`
	}
	if err := ctx.BindJSON(&body); err != nil {
		return
	}
	service.CreateBackup(body.BackupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "create backup success",
		Data: nil,
	})
}

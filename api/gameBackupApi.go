package api

import (
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func DownloadBackup(ctx *gin.Context) {
	service.DownloadBackup(ctx)
}

func GetBackupList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "get backup list success",
		Data: service.GetBackupList(),
	})
}

func RenameBackup(ctx *gin.Context) {

	var body struct {
		FileName string `json:"fileName"`
		NewName  string `json:"newName"`
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

//TODO untest
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
		body.BackupName = ""
	}
	service.CreateBackup(body.BackupName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "create backup success",
		Data: nil,
	})
}

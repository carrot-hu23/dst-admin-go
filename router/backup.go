package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initBackupRouter(router *gin.RouterGroup) {

	backupApi := api.GameBackUpApi{}
	backup := router.Group("/api/game/backup")
	{
		backup.GET("", backupApi.GetBackupList)
		backup.POST("", backupApi.CreateBackup)
		backup.DELETE("", backupApi.DeleteBackup)
		backup.PUT("", backupApi.RenameBackup)
		backup.GET("/download", backupApi.DownloadBackup)
		backup.POST("/upload", backupApi.UploadBackup)

		backup.POST("/snapshot/setting", backupApi.SaveBackupSnapshotsSetting)
		backup.GET("/snapshot/setting", backupApi.GetBackupSnapshotsSetting)
		backup.GET("/snapshot/list", backupApi.BackupSnapshotsList)

	}

}

package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/constant"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/zip"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BackupService2 struct {
	GameConfigService
}

func (this *BackupService2) GetBackupList(ctx *gin.Context) {
	//获取查询参数
	name := ctx.Query("name")
	description := ctx.Query("description")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.DB

	if name, isExist := ctx.GetQuery("name"); isExist {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if description, isExist := ctx.GetQuery("description"); isExist {
		db = db.Where("description LIKE ?", "%"+description+"%")
	}

	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)

	backupList := make([]model.Backup, 0)

	if err := db.Find(&backupList).Error; err != nil {
		fmt.Println(err.Error())
	}

	var total int64
	db2 := database.DB
	if name != "" && description != "" {
		db2.Model(&model.Backup{}).Where("name like ? and description like ?", "%"+name+"%", "%"+description+"%").Count(&total)
	} else if name != "" {
		db2.Model(&model.Backup{}).Where("name like ?", "%"+name+"%").Count(&total)
	} else if description != "" {
		db2.Model(&model.Backup{}).Where("description like ?", "%"+description+"%").Count(&total)
	} else {
		db2.Model(&model.Backup{}).Count(&total)
	}
	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       backupList,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (this *BackupService2) RenameBackup(ctx *gin.Context, fileName, newName string) {

	db := database.DB
	oldBackup := &model.Backup{}
	db.Where("name = ?", fileName).First(oldBackup)
	oldBackup.Name = newName

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup
	err := fileUtils.Rename(path.Join(backupPath, fileName), path.Join(backupPath, newName))
	if err != nil {
		return
	}
	oldBackup.Path = path.Join(backupPath, newName)

	db.Updates(oldBackup)
}

func (this *BackupService2) DeleteBackup(ctx *gin.Context, fileNames []string) {

	db := database.DB
	backupList := make([]model.Backup, 0)
	db.Where("name in ?", fileNames).Find(&backupList)
	for _, backup := range backupList {
		b := model.Backup{}
		db.Where("path = ?", backup.Path).Delete(&b)
		err := fileUtils.DeleteFile(backup.Path)
		if err != nil {
			return
		}
	}
}

// RestoreBackup TODO: 恢复存档
func (this *BackupService2) RestoreBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	filePath := path.Join(cluster.Backup, backupName)
	log.Println("filepath", filePath)

	clusterPath := path.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", cluster.ClusterName)
	err := fileUtils.DeleteDir(clusterPath)
	if err != nil {
		return
	}
	err = zip.Unzip(filePath, clusterPath)
	if err != nil {
		return
	}

}

func (this *BackupService2) CreateBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup

	src := constant.GET_DST_USER_GAME_CONFG_PATH()
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	if backupName == "" {
		gameConfig := vo.NewGameConfigVO()
		this.GetClusterIni(cluster.ClusterName, gameConfig)
		backupName = time.Now().Format("2006-01-02 15:04:05") + "_" + gameConfig.ClusterName + ".zip"
	}
	dst := path.Join(backupPath, backupName)
	log.Println("src", src, dst)
	err := zip.Zip(src, dst)
	if err != nil {
		log.Panicln("create backup error", err)
	}
	log.Println("创建备份成功")
}

func (this *BackupService2) DownloadBackup(c *gin.Context) {
	fileName := c.Query("fileName")

	clusterName := c.GetHeader("level")
	cluster := clusterUtils.GetCluster(clusterName)

	filePath := path.Join(cluster.Backup, fileName)
	//打开文件
	_, err := os.Open(filePath)
	//非空处理
	if err != nil {
		log.Panicln("download filePath error", err)
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	// c.Header("Content-Length", strconv.FormatInt(f.Size(), 10))
	c.File(filePath)
}

func (this *BackupService2) UploadBackup(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	cluster := clusterUtils.GetClusterFromGin(c)
	dst := path.Join(cluster.Backup, file.Filename)

	if fileUtils.Exists(dst) {
		log.Panicln("backup is existed")
	}

	// 上传文件至指定的完整文件路径
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		return
	}

	// c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

}

func (this *BackupService2) backupPath() string {
	dstConfig := dstConfigUtils.GetDstConfig()
	backupPath := dstConfig.Backup
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	return backupPath
}

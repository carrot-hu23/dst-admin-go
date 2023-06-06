package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/zip"
	"dst-admin-go/vo"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

type BackupService struct {
	GameConfigService
}

func (b *BackupService) GetBackupList(ctx *gin.Context) []vo.BackupVo {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	var backupPath = cluster.Backup
	var backupList []vo.BackupVo

	if !fileUtils.Exists(backupPath) {
		return backupList
	}
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(backupPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range fileInfoList {
		if file.IsDir() {
			continue
		}
		suffix := path.Ext(file.Name())
		if suffix == ".zip" || suffix == ".tar" {
			backup := vo.BackupVo{
				FileName:   file.Name(),
				FileSize:   file.Size(),
				CreateTime: file.ModTime(),
				Time:       file.ModTime().Unix(),
			}
			backupList = append(backupList, backup)
		}
	}

	return backupList

}

func (b *BackupService) RenameBackup(ctx *gin.Context, fileName, newName string) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup
	err := fileUtils.Rename(path.Join(backupPath, fileName), path.Join(backupPath, newName))
	if err != nil {
		return
	}
}

func (b *BackupService) DeleteBackup(ctx *gin.Context, fileNames []string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup
	for _, fileName := range fileNames {
		filePath := path.Join(backupPath, fileName)
		if !fileUtils.Exists(filePath) {
			continue
		}
		err := fileUtils.DeleteFile(filePath)
		if err != nil {
			return
		}
	}

}

// TODO: 恢复存档
func (b *BackupService) RestoreBackup(ctx *gin.Context, backupName string) {

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

func (b *BackupService) CreateBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup

	src := constant.GET_DST_USER_GAME_CONFG_PATH()
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	if backupName == "" {
		gameConfig := vo.NewGameConfigVO()
		b.GetClusterIni(cluster.ClusterName, gameConfig)
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

func (b *BackupService) DownloadBackup(c *gin.Context) {
	fileName := c.Query("fileName")

	clusterName := c.GetHeader("cluster")
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

func (b *BackupService) UploadBackup(c *gin.Context) {
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

func (b *BackupService) backupPath() string {
	dstConfig := dstConfigUtils.GetDstConfig()
	backupPath := dstConfig.Backup
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	return backupPath
}

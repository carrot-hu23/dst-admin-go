package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/constant/consts"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/zip"
	"dst-admin-go/vo"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type BackupService struct {
	HomeService
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
		suffix := filepath.Ext(file.Name())
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
	err := fileUtils.Rename(filepath.Join(backupPath, fileName), filepath.Join(backupPath, newName))
	if err != nil {
		return
	}
}

func (b *BackupService) DeleteBackup(ctx *gin.Context, fileNames []string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup
	for _, fileName := range fileNames {
		filePath := filepath.Join(backupPath, fileName)
		if !fileUtils.Exists(filePath) {
			continue
		}
		err := fileUtils.DeleteFile(filePath)
		if err != nil {
			return
		}
	}

}

// RestoreBackup TODO: 恢复存档
func (b *BackupService) RestoreBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	filePath := filepath.Join(cluster.Backup, backupName)
	clusterPath := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", cluster.ClusterName)
	err := fileUtils.DeleteDir(clusterPath)
	if err != nil {
		log.Panicln("删除失败,", clusterPath, err)
	}
	log.Println("正在恢复存档", filePath, filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether"))
	// 先解压到临时目录
	tmpDir := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", "tmp_613e78awhjkdhjkasjkldaso")
	defer func(path string) {
		err := fileUtils.DeleteDir(path)
		if err != nil {
			log.Println("删除tmp失败")
		}
	}(tmpDir)

	err = zip.Unzip(filePath, tmpDir)
	if err != nil {
		log.Panicln("解压失败,", filePath, clusterPath, err)
	}

	tmpFile, err := os.Open(tmpDir)
	if err != nil {
		log.Panicln("打开tmp目录失败,", tmpFile, err)
	}

	var basePath string

	// 遍历文件及其子目录
	err = filepath.Walk(tmpFile.Name(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 找到 Master 目录
		if info.IsDir() && info.Name() == "Master" {
			basePath = filepath.Dir(path)
			return filepath.SkipDir
		}
		return nil
	})
	log.Println(basePath, err)
	if basePath == "" {
		log.Panicln("未找到存档")
	}

	pathList := []string{
		"Master",
		"Caves",
		"cluster.ini",
		"cluster_token.txt",
		"blacklist.txt",
		"adminlist.txt",
	}
	for _, p := range pathList {
		fp := filepath.Join(basePath, p)
		if fileUtils.Exists(fp) {
			fileUtils.Copy(fp, clusterPath)
		}
	}
}

func (b *BackupService) CreateBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup

	src := filepath.Join(consts.KleiDstPath, cluster.ClusterName)
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	if backupName == "" {
		// TODO 增加存档信息
		name := b.GetClusterIni(cluster.ClusterName).ClusterName
		backupName = time.Now().Format("2006-01-02 15:04:05") + "_" + name + ".zip"
	}
	dst := filepath.Join(backupPath, backupName)
	log.Println("src", src, dst)
	err := zip.Zip(src, dst)
	if err != nil {
		log.Panicln("create backup error", err)
	}
	log.Println("创建备份成功")
}

func (b *BackupService) DownloadBackup(c *gin.Context) {
	fileName := c.Query("fileName")

	clusterName := c.GetHeader("world")
	cluster := clusterUtils.GetCluster(clusterName)

	filePath := filepath.Join(cluster.Backup, fileName)
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
	dst := filepath.Join(cluster.Backup, file.Filename)
	log.Println("备份保存在: ", dst)
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

// TODO 备份名称增加存档信息如  猜猜我是谁的世界-10天-spring-1-20-2023071415
func (b *BackupService) GenGameBackUpName(clusterName string) string {
	name := b.GetClusterIni(clusterName).ClusterName
	backupName := time.Now().Format("2006-01-02 15:04:05") + "_" + name + ".zip"

	return backupName
}

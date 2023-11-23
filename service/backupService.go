package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/zip"
	"dst-admin-go/vo"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BackupService struct {
	HomeService
	GameArchive
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

// RestoreBackup TODO: 恢复存档 这里要改
func (b *BackupService) RestoreBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	filePath := filepath.Join(cluster.Backup, backupName)
	clusterPath := filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether", cluster.ClusterName)
	if cluster.ClusterName == "" {
		log.Panicln("cluster is not allow null")
	}
	err := fileUtils.DeleteDir(clusterPath)
	if err != nil {
		log.Panicln("删除失败,", clusterPath, err)
	}
	log.Println("正在恢复存档", filePath, filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether"))

	err = zip.Unzip3(filePath, clusterPath)

	// err = zip.Unzip2(filePath, filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether"), cluster.ClusterName)
	if err != nil {
		log.Panicln("解压失败,", filePath, clusterPath, err)
	}
	// 安装mod
	modoverride, err := fileUtils.ReadFile(dstUtils.GetMasterModoverridesPath(cluster.ClusterName))
	if err != nil {
		log.Println("读取模组失败", err)
	}
	dstUtils.DedicatedServerModsSetup(cluster.ClusterName, modoverride)

}

func (b *BackupService) CreateBackup(ctx *gin.Context, backupName string) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	backupPath := cluster.Backup
	src := dstUtils.GetClusterBasePath(cluster.ClusterName)
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	if backupName == "" {
		backupName = b.GenGameBackUpName(cluster.ClusterName)
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

	clusterName := c.Query("cluster")
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

var SeasonMap = map[string]string{
	"spring": "春天",
	"summer": "夏天",
	"autumn": "秋天",
	"winter": "冬天",
}

// TODO 备份名称增加存档信息如  猜猜我是谁的世界-10天-spring-1-20-2023071415
func (b *BackupService) GenGameBackUpName(clusterName string) string {
	name := b.GetClusterIni(clusterName).ClusterName
	snapshoot := b.Snapshoot(clusterName)

	fmt.Printf("%v\n", snapshoot)

	// 20060102150405_猜猜我是谁的世界_40天_秋季(1/20).zip
	// 猜猜我是谁的房间_季节40天spring(1|20)_模组数量3.zip
	days := strconv.Itoa(snapshoot.Clock.Cycles)
	elapsedDayInSeason := strconv.Itoa(snapshoot.Seasons.ElapsedDaysInSeason)
	seasonDays := strconv.Itoa(snapshoot.Seasons.ElapsedDaysInSeason + snapshoot.Seasons.RemainingDaysInSeason)
	archiveDesc := days + "day_" + SeasonMap[snapshoot.Seasons.Season] + "(" + elapsedDayInSeason + "|" + seasonDays + ")"
	backupName := time.Now().Format("20060102150405") + "_" + name + "_" + archiveDesc + ".zip"

	return backupName
}

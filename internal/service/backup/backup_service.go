package backup

import (
	"dst-admin-go/internal/database"
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/utils/dstUtils"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/pkg/utils/zip"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/game"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type BackupService struct {
	archive     *archive.PathResolver
	dstConfig   dstConfig.Config
	gameProcess game.Process
}

type BackupInfo struct {
	FileName   string    `json:"fileName"`
	FileSize   int64     `json:"fileSize"`
	CreateTime time.Time `json:"createTime"`
	Time       int64     `json:"time"`
}

type BackupSnapshot struct {
	Enable       int `json:"enable"`
	Interval     int `json:"interval"`
	MaxSnapshots int `json:"maxSnapshots"`
	IsCSave      int `json:"isCSave"`
}

func NewBackupService(archive *archive.PathResolver, dstConfig dstConfig.Config, gameProcess game.Process) *BackupService {
	return &BackupService{
		archive:     archive,
		dstConfig:   dstConfig,
		gameProcess: gameProcess,
	}
}

func (b *BackupService) GetBackupList(clusterName string) []BackupInfo {
	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return []BackupInfo{}
	}
	backupPath := config.Backup
	var backupList []BackupInfo

	if !fileUtils.Exists(backupPath) {
		return backupList
	}
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(backupPath)
	if err != nil {
		log.Panicln(err)
	}

	for _, file := range fileInfoList {
		if file.IsDir() {
			continue
		}
		suffix := filepath.Ext(file.Name())
		if suffix == ".zip" || suffix == ".tar" {
			backup := BackupInfo{
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
	clusterName := context.GetClusterName(ctx)
	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}
	backupPath := config.Backup
	err = fileUtils.Rename(filepath.Join(backupPath, fileName), filepath.Join(backupPath, newName))
	if err != nil {
		return
	}
}

func (b *BackupService) DeleteBackup(ctx *gin.Context, fileNames []string) {

	clusterName := context.GetClusterName(ctx)
	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}
	backupPath := config.Backup
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

func (b *BackupService) RestoreBackup(ctx *gin.Context, backupName string) {

	clusterName := context.GetClusterName(ctx)
	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}

	filePath := filepath.Join(config.Backup, backupName)
	clusterPath := filepath.Join(b.archive.ClusterPath(clusterName))
	err = fileUtils.DeleteDir(clusterPath)
	if err != nil {
		log.Panicln("删除失败,", clusterPath, err)
	}
	log.Println("正在恢复存档", filePath, filepath.Join(b.archive.KleiBasePath(clusterName)))

	// err = zip.Unzip2(filePath, filepath.Join(constant.HOME_PATH, ".klei/DoNotStarveTogether"), cluster.ClusterName)
	err = zip.Unzip3(filePath, clusterPath)
	if err != nil {
		log.Panicln("解压失败,", filePath, clusterPath, err)
	}
	// 安装mod
	modoverride, err := fileUtils.ReadFile(b.archive.ModoverridesPath(clusterName, "Master"))
	if err != nil {
		log.Println("读取模组失败", err)
	}
	config, err = b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println(err.Error())
	}
	err = dstUtils.DedicatedServerModsSetup(config, modoverride)
	if err != nil {
		log.Println(err.Error())
	}
}

func (b *BackupService) CreateBackup(clusterName, backupName string) {

	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}
	backupPath := config.Backup

	// 执行 CSave 命令
	err = b.gameProcess.Command(clusterName, "Master", "c_save()")
	if err != nil {
		log.Println("CSave command error:", err)
	}

	// 等待保存完成
	time.Sleep(2 * time.Second)

	src := b.archive.ClusterPath(clusterName)
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	if backupName == "" {
		backupName = b.GenGameBackUpName(clusterName)
	}
	dst := filepath.Join(backupPath, backupName)
	log.Println("src", src, dst)
	err = zip.Zip(src, dst)
	if err != nil {
		log.Panicln("create backup error", err)
	}
	log.Println("创建备份成功")
}

func (b *BackupService) DownloadBackup(c *gin.Context) {
	fileName := c.Query("fileName")

	clusterName := c.GetHeader("level")
	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}

	filePath := filepath.Join(config.Backup, fileName)
	//打开文件
	_, err = os.Open(filePath)
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

	clusterName := context.GetClusterName(c)
	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}
	dst := filepath.Join(config.Backup, file.Filename)

	if fileUtils.Exists(dst) {
		log.Panicln("backup is existed")
	}

	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		return
	}

	// c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

}

func (b *BackupService) ScheduleBackupSnapshots() {
	log.Println("开始创建快照")
	for {
		db := database.Db
		snapshot := model.BackupSnapshot{}
		db.First(&snapshot)
		if snapshot.Enable == 1 {
			if snapshot.Interval == 0 {
				snapshot.Interval = 8
			}
			if snapshot.MaxSnapshots == 0 {
				snapshot.MaxSnapshots = 6
			}
			time.Sleep(time.Duration(snapshot.Interval) * time.Minute)

			// 定时创建备份,每隔 x 分钟备份一次
			// Get the first available cluster config (or "MyCluster" as default)
			clusterName := "MyCluster"
			config, err := b.dstConfig.GetDstConfig(clusterName)
			if err != nil {
				log.Println("failed to get dst config for snapshot:", err)
				continue
			}
			if config.Cluster != "" {
				snapshotPrefix := "(snapshot)"
				if snapshot.IsCSave == 1 {
					// 执行 CSave 命令
					err := b.gameProcess.Command(config.Cluster, "Master", "c_save()")
					if err != nil {
						log.Println("CSave command error:", err)
					}
					// 等待保存完成
					time.Sleep(2 * time.Second)
				}
				b.CreateSnapshotBackup(snapshotPrefix, config.Cluster)
				// 删除快照
				b.DeleteBackupSnapshots(snapshotPrefix, snapshot.MaxSnapshots, config.Cluster, config.Backup)
			}
		} else {
			time.Sleep(1 * time.Minute)
		}
	}
}

func sumMd5(filePath string) string {
	// find save -type f -exec md5sum {} \; | awk '{print $1}' | sort | md5sum
	comamd := "find " + filePath + " -type f -exec md5sum {} \\; | awk '{print $1}' | sort | md5sum"
	info, err := shellUtils.ExecuteCommand(comamd)
	if err != nil {
		return ""
	}
	return info
}

func (b *BackupService) CreateSnapshotBackup(prefix, clusterName string) {

	config, err := b.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		log.Println("failed to get dst config:", err)
		return
	}

	snapshotMd5FilePath := "./snapshotMd5"
	fileUtils.CreateFileIfNotExists(snapshotMd5FilePath)

	src := b.archive.ClusterPath(clusterName)
	dst := filepath.Join(config.Backup, b.GenBackUpSnapshotName(prefix, clusterName))
	log.Println("[Snapshot]正在定时创建游戏备份", "src: ", src, "dst: ", dst)
	err = zip.Zip(src, dst)
	if err != nil {
		log.Println("[Snapshot]create backup error", err)
	}
}

func (b *BackupService) DeleteBackupSnapshots(prefix string, maxSnapshots int, clusterName, backupPath string) {

	log.Println("[Snapshot]正在删除快照备份", "maxSnapshots", maxSnapshots, "clusterName: ", clusterName)

	backupList := b.GetBackupList(clusterName)
	var newBackupList []BackupInfo
	for i := range backupList {
		name := backupList[i].FileName
		if strings.HasPrefix(name, prefix) {
			newBackupList = append(newBackupList, backupList[i])
		}
	}
	if len(newBackupList) > maxSnapshots {
		deleteBackupList := newBackupList[:len(newBackupList)-maxSnapshots]
		for i := range deleteBackupList {
			filePath := filepath.Join(backupPath, deleteBackupList[i].FileName)
			log.Println("删除快照备份", filePath)
			if !fileUtils.Exists(filePath) {
				continue
			}
			err := fileUtils.DeleteFile(filePath)
			if err != nil {
				return
			}
		}
	}

}

func (b *BackupService) backupPath() string {
	// dstConfig := dstConfigUtils.GetDstConfig()
	// backupPath := dstConfig.Backup
	// if !fileUtils.Exists(backupPath) {
	// 	log.Panicln("backup path is not exists")
	// }
	// return backupPath
	return ""
}

var SeasonMap = map[string]string{
	"spring": "春天",
	"summer": "夏天",
	"autumn": "秋天",
	"winter": "冬天",
}

// GenGameBackUpName 备份名称增加存档信息如  猜猜我是谁的世界-10天-spring-1-20-2023071415
func (b *BackupService) GenGameBackUpName(clusterName string) string {
	// 简化实现，使用时间戳和集群名称
	backupName := time.Now().Format("2006年01月02日15点04分05秒") + "_" + clusterName + ".zip"

	return backupName
}

func (b *BackupService) GenBackUpSnapshotName(prefix, clusterName string) string {
	backupName := b.GenGameBackUpName(clusterName)
	backupName = prefix + backupName
	return backupName
}

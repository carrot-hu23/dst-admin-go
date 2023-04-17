package service

import (
	"dst-admin-go/constant"
	archiveutils "dst-admin-go/utils/archiveUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBackupList() []vo.BackupVo {
	var backupPath = dstConfigUtils.GetDstConfig().Backup
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

func RenameBackup(fileName, newName string) {
	backupPath := backupPath()
	fileUtils.Rename(path.Join(backupPath, fileName), path.Join(backupPath, newName))
}

func DeleteBackup(fileNames []string) {
	backupPath := backupPath()
	for _, fileName := range fileNames {
		filePath := path.Join(backupPath, fileName)
		if !fileUtils.Exists(filePath) {
			continue
		}
		fileUtils.DeleteFile(filePath)
	}

}

// TODO: 恢复存档
func RestoreBackup(backupName string) {
	filePath := path.Join(backupPath(), backupName)
	log.Println("filepath", filePath)
	archiveutils.UnZip("C:\\Users\\xm\\Desktop\\backup\\backup\\unzip", filePath)
}

func CreateBackup(backupName string) {
	dstConfig := dstConfigUtils.GetDstConfig()
	backupPath := dstConfig.Backup
	src := constant.GET_DST_USER_GAME_CONFG_PATH()
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	if backupName == "" {
		gameConfig := vo.NewGameConfigVO()
		GetClusterIni(gameConfig)
		backupName = time.Now().Format("2006-01-02 15:04:05") + "_" + gameConfig.ClusterName
	}
	dst := path.Join(backupPath, backupName)
	log.Println("src", src, dst)
	// err := archiveutils.Zip(dst, src)
	// if err != nil {
	// 	log.Panicln(err)
	// }
	archiveutils.Zip2(src, dst)
}

// TODO: 下载存档
func DownloadBackup(c *gin.Context) {
	fileName := c.Query("fileName")

	filePath := path.Join(backupPath(), fileName)
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

func UploadBackup(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	dst := path.Join(backupPath(), file.Filename)

	if fileUtils.Exists(dst) {
		log.Panicln("backup is existed")
	}

	// 上传文件至指定的完整文件路径
	c.SaveUploadedFile(file, dst)

	// c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

}

func backupPath() string {
	dstConfig := dstConfigUtils.GetDstConfig()
	backupPath := dstConfig.Backup
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	return backupPath
}

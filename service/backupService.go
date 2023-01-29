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

	"github.com/gin-gonic/gin"
)

func GetBackupList() []vo.BackupVo {

	log.Println("cluster", constant.GET_DST_USER_GAME_CONFG_PATH())

	var backupPath = dstConfigUtils.GetDstConfig().Backup
	var backupList []vo.BackupVo

	log.Println("backup path", backupPath)

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

// func checkBackupPath(fileName string) string {
// 	backupPath := dstConfigUtils.GetDstConfig().Backup
// 	if fileUtils.Exists(backupPath) {
// 		return backupPath + fileName
// 	}
// 	return ""
// }

func RenameBackup(fileName, newName string) {
	path := backupPath()
	log.Println(path, fileName)
	fileUtils.Rename(path+fileName, path+newName)
}

func DeleteBackup(fileNames []string) {
	backupPath := backupPath()
	for _, fileName := range fileNames {
		path := backupPath + fileName
		if !fileUtils.Exists(path) {
			continue
		}
		fileUtils.DeleteFile(path)
	}

}

// TODO: 恢复存档
func RestoreBackup(backupName string) {
	filePath := backupPath() + backupName
	log.Println("filepath", filePath)
	archiveutils.UnZip("C:\\Users\\xm\\Desktop\\backup\\backup\\unzip", filePath)
}

func CreateBackup(backupName string) {
	dstConfig := dstConfigUtils.GetDstConfig()
	backupPath := dstConfig.Backup
	src := dstConfig.DoNotStarveTogether + dstConfig.Cluster
	if !fileUtils.Exists(backupPath) {
		log.Panicln("backup path is not exists")
	}
	dst := backupPath + backupName
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

	filePath := backupPath() + fileName
	//打开文件
	_, err := os.Open(filePath)
	//非空处理
	if err != nil {
		log.Panicln("download filePath error", err)
	}
	// f, e := os.Stat(filePath)
	// if e != nil {
	// 	log.Println(e)
	// }
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

	dst := backupPath() + file.Filename

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

package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DstMapApi struct {
}

func (d *DstMapApi) GenDstMap(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	outputImage := filepath.Join(dstUtils.GetClusterBasePath(cluster.ClusterName), "dst_map.png")
	sessionPath := filepath.Join(dstUtils.GetKleiDstPath(), cluster.ClusterName, "Master", "save", "session")
	filePath, err := findLatestMetaFile(sessionPath)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("生成地图", filePath, outputImage)
	generator := service.NewDSTMapGenerator()
	height, width, err := service.ExtractDimensions(filePath)
	if err != nil {
		log.Panicln(err)
	}
	err = generator.GenerateMap(
		filePath,
		outputImage,
		height,
		width,
	)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (d *DstMapApi) GetDstMapImage(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	outputImage := filepath.Join(dstUtils.GetClusterBasePath(cluster.ClusterName), "dst_map.png")
	log.Println(outputImage)
	// 使用 Gin 提供的文件传输方法返回图片
	ctx.File(outputImage)
	ctx.Header("Content-Type", "image/png")
}

func (d *DstMapApi) HasWalrusHutPlains(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	sessionPath := filepath.Join(dstUtils.GetKleiDstPath(), cluster.ClusterName, "Master", "save", "session")
	filePath, err := findLatestMetaFile(sessionPath)
	if err != nil {
		log.Panicln(err)
	}
	file, err := fileUtils.ReadFile(filePath)
	if err != nil {
		log.Panicln(err)
	}
	hasWalrusHutPlains := strings.Contains(file, "WalrusHut_Plains")
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: hasWalrusHutPlains,
	})

}

func findLatestMetaFile(directory string) (string, error) {
	// 检查指定目录是否存在
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("目录不存在：%s", directory)
	}

	// 获取指定目录下一级的所有子目录
	subdirs, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("读取目录失败：%s", err)
	}

	// 用于存储最新的.meta文件路径和其修改时间
	var latestMetaFile string
	var latestMetaFileTime time.Time

	for _, subdir := range subdirs {
		// 检查子目录是否是目录
		if subdir.IsDir() {
			subdirPath := filepath.Join(directory, subdir.Name())

			// 获取子目录下的所有文件
			files, err := ioutil.ReadDir(subdirPath)
			if err != nil {
				return "", fmt.Errorf("读取子目录失败：%s", err)
			}

			for _, file := range files {
				// 检查文件是否是.meta文件
				if !file.IsDir() && filepath.Ext(file.Name()) != ".meta" {
					// 获取文件的修改时间
					modifiedTime := file.ModTime()

					// 如果找到的文件的修改时间比当前最新的.meta文件的修改时间更晚，则更新最新的.meta文件路径和修改时间
					if modifiedTime.After(latestMetaFileTime) {
						latestMetaFile = filepath.Join(subdirPath, file.Name())
						latestMetaFileTime = modifiedTime
					}
				}
			}
		}
	}

	if latestMetaFile == "" {
		return "", fmt.Errorf("未找到文件")
	}

	return latestMetaFile, nil
}

package handler

import (
	"dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/dstMap"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type DstMapHandler struct {
	generator       *dstMap.DSTMapGenerator
	archiveResolver *archive.PathResolver
}

func NewDstMapHandler(archiveResolver *archive.PathResolver, generator *dstMap.DSTMapGenerator) *DstMapHandler {
	return &DstMapHandler{
		archiveResolver: archiveResolver,
		generator:       generator,
	}
}

func (d *DstMapHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/dst/map/gen", d.GenDstMap)
	router.GET("/api/dst/map/image", d.GetDstMapImage)
	router.GET("/api/dst/map/has/walrusHut/plains", d.HasWalrusHutPlains)
	router.GET("/api/dst/map/session/file", d.GetSessionFile)
	router.GET("/api/dst/map/player/session/file", d.GetPlayerSessionFile)
}

// GenDstMap 生成地图 生成 swagger 文档注释
// @Summary 生成地图
// @Description 生成地图
// @Tags dstMap
// @Param levelName query string true "levelName"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/dst/map/gen [get]
func (d *DstMapHandler) GenDstMap(ctx *gin.Context) {

	levelName := ctx.Query("levelName")
	if levelName == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "levelName 参数不能为空",
		})
		return
	}
	clusterName := context.GetClusterName(ctx)
	clusterPath := d.archiveResolver.ClusterPath(clusterName)
	outputImage := filepath.Join(clusterPath, "dst_map_"+levelName+".jpg")
	sessionPath := filepath.Join(clusterPath, levelName, "save", "session")
	filePath, err := findLatestMetaFile(sessionPath)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("生成地图", filePath, outputImage)
	height, width, err := dstMap.ExtractDimensions(filePath)
	if err != nil {
		log.Panicln(err)
	}
	err = d.generator.GenerateMap(
		filePath,
		outputImage,
		height,
		width,
	)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// GetDstMapImage 获取地图图片 获取 swagger 文档注释
// @Summary 获取地图图片
// @Description 获取地图图片
// @Tags dstMap
// @Param levelName query string true "levelName"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/dst/map/image [get]
func (d *DstMapHandler) GetDstMapImage(ctx *gin.Context) {

	levelName := ctx.Query("levelName")
	if levelName == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "levelName 参数不能为空",
		})
		return
	}

	clusterName := context.GetClusterName(ctx)
	clusterPath := d.archiveResolver.ClusterPath(clusterName)

	outputImage := filepath.Join(clusterPath, "dst_map_"+levelName+".jpg")
	log.Println(outputImage)
	// 使用 Gin 提供的文件传输方法返回图片
	ctx.File(outputImage)
	ctx.Header("Content-Type", "image/png")

}

// HasWalrusHutPlains 检测地图中是否有walrusHutPlains 获取 swagger 文档注释
// @Summary 检测地图中是否有walrusHutPlains
// @Description 检测地图中是否有walrusHutPlains
// @Tags dstMap
// @Param levelName query string true "levelName"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/dst/map/has/walrusHut/plains [get]
func (d *DstMapHandler) HasWalrusHutPlains(ctx *gin.Context) {

	levelName := ctx.Query("levelName")
	if levelName == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "levelName 参数不能为空",
		})
		return
	}

	clusterName := context.GetClusterName(ctx)
	clusterPath := d.archiveResolver.ClusterPath(clusterName)

	sessionPath := filepath.Join(clusterPath, levelName, "save", "session")
	filePath, err := findLatestMetaFile(sessionPath)
	if err != nil {
		log.Panicln(err)
	}
	file, err := fileUtils.ReadFile(filePath)
	if err != nil {
		log.Panicln(err)
	}
	hasWalrusHutPlains := strings.Contains(file, "WalrusHut_Plains")
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: hasWalrusHutPlains,
	})

}

// GetSessionFile 获取存档文件 获取 swagger 文档注释
// @Summary 获取存档文件
// @Description 获取存档文件
// @Tags dstMap
// @Param levelName query string true "levelName"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/dst/map/session/file [get]
func (d *DstMapHandler) GetSessionFile(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)
	clusterPath := d.archiveResolver.ClusterPath(clusterName)

	levelName := ctx.Query("levelName")
	if levelName == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "levelName 参数不能为空",
		})
		return
	}
	sessionPath := filepath.Join(clusterPath, levelName, "save", "session")
	filePath, err := findLatestMetaFile(sessionPath)
	if err != nil {
		log.Panicln(err)
	}
	file, err := fileUtils.ReadFile(filePath)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: file,
	})

}

// GetPlayerSessionFile 获取玩家存档文件 获取 swagger 文档
// @Summary 获取玩家存档文件
// @Description 获取玩家存档文件
// @Tags dstMap
// @Param levelName query string true "levelName"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/dst/map/player/session/file [get]
func (d *DstMapHandler) GetPlayerSessionFile(ctx *gin.Context) {

	clusterName := context.GetClusterName(ctx)
	clusterPath := d.archiveResolver.ClusterPath(clusterName)

	levelName := ctx.Query("levelName")
	kuId := ctx.Query("kuId")
	if levelName == "" || kuId == "" {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  "levelName or kuId 参数不能为空",
		})
		return
	}

	baseSessionFile := filepath.Join(clusterPath, levelName, "save", "session")
	latestMetaFile, err2 := findLatestMetaFile(baseSessionFile)
	if err2 != nil {
		ctx.JSON(http.StatusBadRequest, response.Response{
			Code: 400,
			Msg:  err2.Error(),
		})
		return
	}
	sessionID := extractSessionID(latestMetaFile)
	sessionPath := filepath.Join(baseSessionFile, sessionID, kuId+"_")
	log.Println(sessionPath)
	filePath, err := findLatestPlayerFile(sessionPath)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(filePath)
	file, err := fileUtils.ReadFile(filePath)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: file,
	})

}

func findLatestPlayerFile(directory string) (string, error) {
	// 检查指定目录是否存在
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("目录不存在：%s", directory)
	}

	// 用于存储最新的.meta文件路径和其修改时间
	var latestFile string
	var latestFileTime time.Time

	// 获取指定目录下所有的文件
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("读取目录失败：%s", err)
	}
	for _, file := range files {
		// 检查文件是否是文件
		if !file.IsDir() {
			// 获取文件的修改时间
			modifiedTime := file.ModTime()
			// 如果找到的文件的修改时间比当前最新的.meta文件的修改时间更晚，则更新最新的.meta文件路径和修改时间
			if modifiedTime.After(latestFileTime) {
				latestFile = filepath.Join(directory, file.Name())
				latestFileTime = modifiedTime
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("未找到文件")
	}

	return latestFile, nil
}

func extractSessionPrefix(sessionFile string) string {
	parts := strings.Split(sessionFile, "/")
	if len(parts) >= 2 {
		return parts[0] + "/" + parts[1]
	}
	return sessionFile
}

const sessionPrefix = "/save/session/"

// extractSessionID 提取 /save/session/ 后的第一个路径段（如 925F2AFB73839B9E）
func extractSessionID(p string) string {
	// 找到 "/save/session/" 的起始位置
	i := strings.Index(p, sessionPrefix)
	// 由于题目保证一定存在，可直接跳过错误检查
	rest := p[i+len(sessionPrefix):]
	// 取第一个 '/' 之前的部分（即 session ID）
	if j := strings.Index(rest, "/"); j != -1 {
		return rest[:j]
	}
	// 理论上不会走到这里（因为后面还有子目录如 /0000000002），但为安全起见：
	return rest // 整个剩余部分（如路径恰好以 ID 结尾）
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

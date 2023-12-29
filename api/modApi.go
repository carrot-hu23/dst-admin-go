package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/mod"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ModApi struct {
}

func (m *ModApi) SearchModList(ctx *gin.Context) {

	//获取查询参数
	text := ctx.Query("text")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	data, err := mod.SearchModList(text, page, size)
	if err != nil {
		log.Panicln("搜索mod失败", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

func (m *ModApi) GetModInfo(ctx *gin.Context) {

	moId := ctx.Param("modId")
	modinfo, err, status := mod.SubscribeModByModId(moId)
	if err != nil {
		log.Panicln("模组下载失败", "status: ", status)
	}
	var mod_config map[string]interface{}
	_ = json.Unmarshal([]byte(modinfo.ModConfig), &mod_config)
	modData := map[string]interface{}{
		"auth":          modinfo.Auth,
		"consumer_id":   modinfo.ConsumerAppid,
		"creator_appid": modinfo.CreatorAppid,
		"description":   modinfo.Description,
		"file_url":      modinfo.FileUrl,
		"modid":         modinfo.Modid,
		"img":           modinfo.Img,
		"last_time":     modinfo.LastTime,
		"name":          modinfo.Name,
		"v":             modinfo.V,
		"mod_config":    mod_config,
		"update":        modinfo.Update,
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: modData,
	})
}

func (m *ModApi) GetMyModList(ctx *gin.Context) {

	var modInfos []model.ModInfo
	db := database.DB

	db.Find(&modInfos)

	var modDataList []map[string]interface{}

	//var workshopIds []string
	//for i := range modInfos {
	//	workshopIds = append(workshopIds, modInfos[i].Modid)
	//}
	//publishedFileDetails, err := mod.GetPublishedFileDetailsWithGet(workshopIds)

	for _, modinfo := range modInfos {

		//update := false
		//tags := []string{}
		//if err == nil {
		//	for i := range publishedFileDetails {
		//		publishedfiledetail := publishedFileDetails[i]
		//		if modinfo.Modid == publishedfiledetail.Publishedfileid && modinfo.LastTime < publishedfiledetail.TimeUpdated {
		//			update = true
		//			for _, tag := range publishedfiledetail.Tags {
		//				tags = append(tags, tag.Tag)
		//			}
		//		}
		//	}
		//}

		var mod_config map[string]interface{}
		_ = json.Unmarshal([]byte(modinfo.ModConfig), &mod_config)

		modData := map[string]interface{}{
			"auth":          modinfo.Auth,
			"consumer_id":   modinfo.ConsumerAppid,
			"creator_appid": modinfo.CreatorAppid,
			"description":   modinfo.Description,
			"file_url":      modinfo.FileUrl,
			"modid":         modinfo.Modid,
			"img":           modinfo.Img,
			"last_time":     modinfo.LastTime,
			"name":          modinfo.Name,
			"v":             modinfo.V,
			"mod_config":    mod_config,
			"update":        modinfo.Update,
		}
		modDataList = append(modDataList, modData)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: modDataList,
	})

}

func (m *ModApi) UpdateAllModInfos(ctx *gin.Context) {

	mod.UpdateModinfoList()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *ModApi) DeleteMod(ctx *gin.Context) {

	modId := ctx.Param("modId")
	db := database.DB
	db.Where("modid = ?", modId).Delete(&model.ModInfo{})

	dstConfig := dstConfigUtils.GetDstConfig()
	mod_download_path := dstConfig.Mod_download_path
	mod_path := filepath.Join(mod_download_path, "/steamapps/workshop/content/322330/", modId)
	fileUtils.DeleteDir(mod_path)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: modId,
	})
}

func (m *ModApi) DeleteSetupWorkshop(ctx *gin.Context) {
	dstPath := dstConfigUtils.GetDstConfig().Force_install_dir
	modsPath := filepath.Join(dstPath, "mods")
	// 删除所有workshop-xxx mod

	directories, err := fileUtils.ListDirectories(modsPath)
	if err != nil {
		log.Panicln("delete dst workshop file error", err)
	}
	var workshopList []string
	for _, directory := range directories {
		if strings.Contains(directory, "workshop") {
			workshopList = append(workshopList, directory)
		}
	}
	for _, workshop := range workshopList {
		err := fileUtils.DeleteDir(workshop)
		if err != nil {
			return
		}
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (m *ModApi) GetModInfoFile(ctx *gin.Context) {
	modId := ctx.Param("modId")
	var modInfos model.ModInfo
	db := database.DB
	db.Where("modid = ?", modId).Find(&modInfos)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: modInfos,
	})
}

func (m *ModApi) SaveModInfoFile(ctx *gin.Context) {

	var modInfos model.ModInfo
	err := ctx.ShouldBind(&modInfos)
	if err != nil {
		log.Panicln("参数解析失败")
	}
	db := database.DB
	db.Save(&modInfos)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: modInfos,
	})
}

func (m *ModApi) UpdateMod(ctx *gin.Context) {

	modId := ctx.Param("modId")
	db := database.DB
	db.Where("modid = ?", modId).Delete(&model.ModInfo{})

	dstConfig := dstConfigUtils.GetDstConfig()
	mod_download_path := dstConfig.Mod_download_path
	mod_path := filepath.Join(mod_download_path, "/steamapps/workshop/content/322330/", modId)
	fileUtils.DeleteDir(mod_path)

	modinfo, err, status := mod.SubscribeModByModId(modId)
	if err != nil {
		log.Panicln("模组下载失败", "status: ", status)
	}
	var mod_config map[string]interface{}
	json.Unmarshal([]byte(modinfo.ModConfig), &mod_config)
	mod := map[string]interface{}{
		"auth":          modinfo.Auth,
		"consumer_id":   modinfo.ConsumerAppid,
		"creator_appid": modinfo.CreatorAppid,
		"description":   modinfo.Description,
		"file_url":      modinfo.FileUrl,
		"modid":         modinfo.Modid,
		"img":           modinfo.Img,
		"last_time":     modinfo.LastTime,
		"name":          modinfo.Name,
		"v":             modinfo.V,
		"mod_config":    mod_config,
		"update":        modinfo.Update,
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: mod,
	})
}

// AddModInfoFile 手动添加模组
func (m *ModApi) AddModInfoFile(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)

	var payload struct {
		WorkshopId string `json:"workshopId"`
		Modinfo    string `json:"modinfo"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln("参数解析失败")
	}
	if payload.WorkshopId == "" {
		log.Panicln("workshopId can not be null")
	}

	// 创建workshop文件
	workshopDirPath := filepath.Join(cluster.ModDownloadPath, "/steamapps/workshop/content/322330", payload.WorkshopId)
	fileUtils.CreateDirIfNotExists(workshopDirPath)
	modinfoPath := filepath.Join(workshopDirPath, "modinfo.lua")

	err = fileUtils.CreateFileIfNotExists(modinfoPath)
	if err != nil {
		log.Panicln("创建 modinfo.lua 失败", modinfoPath, err)
	}
	err = fileUtils.WriterTXT(modinfoPath, payload.Modinfo)
	if err != nil {
		log.Panicln("写入 modinfo.lua 失败， path: ", modinfoPath, "error: ", err)
	}

	mod.AddModInfo(payload.WorkshopId)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

// AddModInfoFile 手动添加模组
func (m *ModApi) UploadModFile(ctx *gin.Context) {

}

package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/mod"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	mod_path := filepath.Join(mod_download_path, "steamapps", "workshop", "content", "322330", modId)
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

const (
	steamAPIKey = "73DF9F781D195DFD3D19DED1CB72EEE6"
)

type WorkshopItemDetail struct {
	WorkShopId  string  `json:"workshopId"`
	Name        string  `json:"name"`
	Timeupdated int64   `json:"timeupdated"`
	Timelast    float64 `json:"timelast"`
	Img         string  `json:"img"`
}

type AppworkshopAcf struct {
}

func (m *ModApi) GetUgcModAcf(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	levelName := ctx.Query("levelName")

	acfPath := dstUtils.GetUgcAcfPath(cluster.ClusterName, levelName)
	log.Println("acfPath", acfPath)
	acfWorkshops := dstUtils.ParseACFFile(acfPath)

	var workshopItemDetails []WorkshopItemDetail

	var modIds []string
	for key := range acfWorkshops {
		modIds = append(modIds, key)
	}

	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	for i := range modIds {
		data.Set("publishedfileids["+strconv.Itoa(i)+"]", modIds[i])
	}
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Panicln(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Panicln(err)
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok {
		log.Panicln(err)
	}
	for i := range dataList {
		workshop := dataList[i].(map[string]interface{})
		_, find := workshop["time_updated"]
		if find {
			timeUpdated := workshop["time_updated"].(float64)
			modId := workshop["publishedfileid"].(string)
			value, ok := acfWorkshops[modId]
			if ok {
				img := workshop["preview_url"].(string)
				img = fmt.Sprintf("%s?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%%23000000&letterbox=true", img)
				workshopItemDetails = append(workshopItemDetails, WorkshopItemDetail{
					WorkShopId:  modId,
					Timeupdated: value.TimeUpdated,
					Timelast:    timeUpdated,
					Img:         img,
					Name:        workshop["title"].(string),
				})
			}
		}
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: workshopItemDetails,
	})
}

func (m *ModApi) DeleteUgcModFile(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")
	workshopId := ctx.Query("workshopId")

	modFilePath := dstUtils.GetUgcWorkshopModPath(clusterName, levelName, workshopId)

	log.Println("modFilePath", modFilePath)
	if fileUtils.Exists(modFilePath) {
		err := fileUtils.DeleteDir(modFilePath)
		if err != nil {
			log.Panicln(err)
		}
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

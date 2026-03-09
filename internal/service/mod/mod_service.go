package mod

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/pkg/utils/shellUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/dstConfig"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
)

const (
	steamAPIKey = "73DF9F781D195DFD3D19DED1CB72EEE6"
	appID       = 322330
	language    = 6
)

type ModService struct {
	db           *gorm.DB
	dstConfig    dstConfig.Config
	pathResolver *archive.PathResolver
}

func NewModService(db *gorm.DB, config dstConfig.Config, pathResolver *archive.PathResolver) *ModService {
	return &ModService{
		db:           db,
		dstConfig:    config,
		pathResolver: pathResolver,
	}
}

// SearchResult 搜索结果
type SearchResult struct {
	Page      int       `json:"page"`
	Size      int       `json:"size"`
	Total     int       `json:"total"`
	TotalPage int       `json:"totalPage"`
	Data      []ModInfo `json:"data"`
}

// ModInfo 搜索返回的mod信息
type ModInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Author        string  `json:"author"`
	Desc          string  `json:"desc"`
	Time          int     `json:"time"`
	Sub           int     `json:"sub"`
	Img           string  `json:"img"`
	FileUrl       string  `json:"file_url"`
	V             string  `json:"v"`
	LastTime      float64 `json:"last_time"`
	ConsumerAppid float64 `json:"consumer_appid"`
	CreatorAppid  float64 `json:"creator_appid"`
	Vote          struct {
		Star int `json:"star"`
		Num  int `json:"num"`
	} `json:"vote"`
	Child []string `json:"child,omitempty"`
}

// Publishedfiledetail Steam API 返回的模组详情
type Publishedfiledetail struct {
	Publishedfileid       string  `json:"publishedfileid"`
	Result                int     `json:"result"`
	Creator               string  `json:"creator"`
	CreatorAppID          int     `json:"creator_app_id"`
	ConsumerAppID         int     `json:"consumer_app_id"`
	Filename              string  `json:"filename"`
	FileURL               string  `json:"file_url"`
	HcontentFile          string  `json:"hcontent_file"`
	PreviewURL            string  `json:"preview_url"`
	HcontentPreview       string  `json:"hcontent_preview"`
	Title                 string  `json:"title"`
	Description           string  `json:"description"`
	TimeCreated           float64 `json:"time_created"`
	TimeUpdated           float64 `json:"time_updated"`
	Visibility            int     `json:"visibility"`
	BanReason             string  `json:"ban_reason"`
	Subscriptions         int     `json:"subscriptions"`
	Favorited             int     `json:"favorited"`
	LifetimeSubscriptions int     `json:"lifetime_subscriptions"`
	LifetimeFavorited     int     `json:"lifetime_favorited"`
	Views                 int     `json:"views"`
	Tags                  []struct {
		Tag string `json:"tag"`
	} `json:"tags"`
}

// WorkshopItemDetail UGC mod详情
type WorkshopItemDetail struct {
	WorkShopId  string  `json:"workshopId"`
	Name        string  `json:"name"`
	Timeupdated int64   `json:"timeupdated"`
	Timelast    float64 `json:"timelast"`
	Img         string  `json:"img"`
}

// WorkshopItem ACF文件中的Workshop项
type WorkshopItem struct {
	TimeUpdated int64
	Manifest    string
	Ugchandle   string
}

// SearchModList 搜索模组列表
func (s *ModService) SearchModList(text string, page, size int, lang string) (*SearchResult, error) {
	// 判断是否是modID搜索
	modId, ok := isModId(text)
	if ok {
		modInfo := s.searchModInfoByWorkshopId(modId)
		data := []ModInfo{}
		if modInfo.ID != "" {
			data = append(data, modInfo)
		}
		return &SearchResult{
			Page:      1,
			Size:      1,
			Total:     1,
			TotalPage: 1,
			Data:      data,
		}, nil
	}

	// 调用 Steam API 搜索
	urlStr := "http://api.steampowered.com/IPublishedFileService/QueryFiles/v1/"
	data := url.Values{
		"page":             {fmt.Sprintf("%d", page)},
		"key":              {steamAPIKey},
		"appid":            {"322330"},
		"language":         {"6"},
		"return_tags":      {"true"},
		"numperpage":       {fmt.Sprintf("%d", size)},
		"search_text":      {text},
		"return_vote_data": {"true"},
		"return_children":  {"true"},
	}
	if lang == "zh" {
		data.Set("language", "6")
	} else {
		data.Set("language", "")
	}
	urlStr = urlStr + "?" + data.Encode()

	var modData map[string]interface{}
	for i := 0; i < 2; i++ {
		resp, err := http.Get(urlStr)
		if err != nil {
			return nil, fmt.Errorf("搜索mod失败: %w", err)
		}
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&modData)
		if err != nil {
			return nil, fmt.Errorf("解析mod数据失败: %w", err)
		}
		if modData["response"] != nil {
			break
		}
	}

	if modData["response"] == nil {
		return nil, errors.New("no response found in mod data")
	}

	modResponse := modData["response"].(map[string]interface{})
	total := int(modResponse["total"].(float64))
	modInfoRaw := modResponse["publishedfiledetails"].([]interface{})

	modList := make([]ModInfo, 0)
	if len(modInfoRaw) > 0 {
		for _, modInfoRaw := range modInfoRaw {
			modInfo := modInfoRaw.(map[string]interface{})
			img := modInfo["preview_url"].(string)
			voteData := modInfo["vote_data"].(map[string]interface{})
			auth := modInfo["creator"].(string)
			var authorURL string
			if auth != "" {
				authorURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s/?xml=1", auth)
			}
			mod := ModInfo{
				ID:     fmt.Sprintf("%v", modInfo["publishedfileid"]),
				Name:   fmt.Sprintf("%v", modInfo["title"]),
				Author: authorURL,
				Desc:   fmt.Sprintf("%v", modInfo["file_description"]),
				Time:   int(modInfo["time_updated"].(float64)),
				Sub:    int(modInfo["subscriptions"].(float64)),
				Img:    img,
				Vote: struct {
					Star int `json:"star"`
					Num  int `json:"num"`
				}{
					Star: int(voteData["score"].(float64)*5) + 1,
					Num:  int(voteData["votes_up"].(float64) + voteData["votes_down"].(float64)),
				},
			}
			if modInfo["num_children"].(float64) != 0 {
				children := modInfo["children"].([]interface{})
				child := make([]string, len(children))
				for i, c := range children {
					child[i] = fmt.Sprintf("%v", c.(map[string]interface{})["publishedfileid"])
				}
				mod.Child = child
			}
			modList = append(modList, mod)
		}
	}

	return &SearchResult{
		Page:      page,
		Size:      size,
		Total:     total,
		TotalPage: int(math.Ceil(float64(total) / float64(size))),
		Data:      modList,
	}, nil
}

// SubscribeModByModId 订阅并下载模组
func (s *ModService) SubscribeModByModId(clusterName, modId, lang string) (*model.ModInfo, error) {
	if !isWorkshopId(modId) {
		// 非workshop mod，从本地读取
		return s.getLocalModInfo(clusterName, lang, modId)
	}

	// 从Steam API获取mod信息
	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	data.Set("publishedfileids[0]", modId)
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求Steam API失败: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok || len(dataList) == 0 {
		return nil, errors.New("获取mod信息失败")
	}

	data2 := dataList[0].(map[string]interface{})
	img := data2["preview_url"].(string)
	auth := data2["creator"].(string)
	var authorURL string
	if auth != "" {
		authorURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s/?xml=1", auth)
	}

	name := data2["title"].(string)
	lastTime := data2["time_updated"].(float64)
	description := data2["file_description"].(string)
	auth = authorURL
	fileUrl := data2["file_url"]
	img = fmt.Sprintf("%s?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%%23000000&letterbox=true", img)
	v := s.getVersion(data2["tags"])
	creatorAppid := data2["creator_appid"].(float64)
	consumerAppid := data2["consumer_appid"].(float64)

	// 检查数据库中是否已存在
	existingMod, err := s.GetModByModId(modId)
	if err == nil && existingMod.Modid != "" {
		if lastTime == existingMod.LastTime {
			return existingMod, nil
		}
		// 需要更新
		var modConfig string
		var fileUrlStr = ""
		if fileUrl != nil {
			fileUrlStr = fileUrl.(string)
		}
		if fileUrlStr != "" {
			modConfigJson, _ := json.Marshal(s.getV1ModInfoConfig(clusterName, lang, modId, fileUrlStr))
			modConfig = string(modConfigJson)
		} else {
			modConfigJson, _ := json.Marshal(s.getModInfoConfig(clusterName, lang, modId))
			modConfig = string(modConfigJson)
		}

		existingMod.LastTime = lastTime
		existingMod.Name = name
		existingMod.Auth = auth
		existingMod.Description = description
		existingMod.Img = img
		existingMod.V = v
		existingMod.ModConfig = modConfig
		existingMod.Update = false
		s.db.Save(existingMod)
		return existingMod, nil
	}

	// 新增mod
	var fileUrlStr = ""
	if fileUrl != nil {
		fileUrlStr = fileUrl.(string)
	}

	var modConfig string
	if fileUrlStr != "" {
		modConfigJson, _ := json.Marshal(s.getV1ModInfoConfig(clusterName, lang, modId, fileUrlStr))
		modConfig = string(modConfigJson)
	} else {
		modConfigJson, _ := json.Marshal(s.getModInfoConfig(clusterName, lang, modId))
		modConfig = string(modConfigJson)
	}

	newModInfo := &model.ModInfo{
		Auth:          auth,
		ConsumerAppid: consumerAppid,
		CreatorAppid:  creatorAppid,
		Description:   description,
		FileUrl:       fileUrlStr,
		Modid:         modId,
		Img:           img,
		LastTime:      lastTime,
		Name:          name,
		V:             v,
		ModConfig:     modConfig,
	}

	err = s.db.Create(newModInfo).Error
	return newModInfo, err
}

// GetMyModList 获取已订阅的模组列表
func (s *ModService) GetMyModList() ([]model.ModInfo, error) {
	var modInfos []model.ModInfo
	err := s.db.Find(&modInfos).Error
	return modInfos, err
}

// GetModByModId 根据modId获取模组
func (s *ModService) GetModByModId(modId string) (*model.ModInfo, error) {
	var modInfo model.ModInfo
	err := s.db.Where("modid = ?", modId).First(&modInfo).Error
	return &modInfo, err
}

// DeleteMod 删除模组
func (s *ModService) DeleteMod(clusterName, modId string) error {
	// 从数据库删除
	err := s.db.Where("modid = ?", modId).Delete(&model.ModInfo{}).Error
	if err != nil {
		return err
	}

	// 删除本地文件
	config, _ := s.dstConfig.GetDstConfig(clusterName)
	modDownloadPath := config.Mod_download_path
	modPath := filepath.Join(modDownloadPath, "steamapps", "workshop", "content", "322330", modId)
	return fileUtils.DeleteDir(modPath)
}

// UpdateAllModInfos 批量更新所有模组信息
func (s *ModService) UpdateAllModInfos(clusterName, lang string) error {
	var modInfos []model.ModInfo
	var needUpdateList []model.ModInfo
	var workshopIds []string

	s.db.Find(&modInfos)

	for i := range modInfos {
		workshopIds = append(workshopIds, modInfos[i].Modid)
	}

	publishedFileDetails, err := s.getPublishedFileDetailsBatched(workshopIds, 20)
	if err != nil {
		return err
	}

	for i := range publishedFileDetails {
		publishedfiledetail := publishedFileDetails[i]
		for j := range modInfos {
			if modInfos[j].Modid == publishedfiledetail.Publishedfileid && modInfos[j].LastTime < publishedfiledetail.TimeUpdated {
				needUpdateList = append(needUpdateList, modInfos[i])
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(needUpdateList))

	for i := range needUpdateList {
		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
				wg.Done()
			}()
			modId := needUpdateList[i].Modid
			// 删除之前的数据
			config, _ := s.dstConfig.GetDstConfig(clusterName)
			modDownloadPath := config.Mod_download_path
			modPath := filepath.Join(modDownloadPath, "/steamapps/workshop/content/322330/", modId)
			_ = fileUtils.DeleteDir(modPath)
			_, _ = s.SubscribeModByModId(clusterName, modId, lang)
		}(i)
	}
	wg.Wait()

	return nil
}

// DeleteSetupWorkshop 删除所有workshop模组
func (s *ModService) DeleteSetupWorkshop(clusterName string) error {
	config, err := s.dstConfig.GetDstConfig(clusterName)
	if err != nil {
		return err
	}

	dstPath := config.Force_install_dir
	modsPath := filepath.Join(dstPath, "mods")

	directories, err := fileUtils.ListDirectories(modsPath)
	if err != nil {
		return fmt.Errorf("列出目录失败: %w", err)
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
			return err
		}
	}

	return nil
}

// SaveModInfo 保存模组信息
func (s *ModService) SaveModInfo(modInfo *model.ModInfo) error {
	return s.db.Save(modInfo).Error
}

// AddModInfo 手动添加模组
func (s *ModService) AddModInfo(clusterName, lang, modid, modinfo, modDownloadPath string) error {
	// 创建workshop文件
	workshopDirPath := filepath.Join(modDownloadPath, "/steamapps/workshop/content/322330", modid)
	fileUtils.CreateDirIfNotExists(workshopDirPath)

	modinfoPath := filepath.Join(workshopDirPath, "modinfo.lua")
	err := fileUtils.CreateFileIfNotExists(modinfoPath)
	if err != nil {
		return fmt.Errorf("创建modinfo.lua失败: %w", err)
	}

	err = fileUtils.WriterTXT(modinfoPath, modinfo)
	if err != nil {
		return fmt.Errorf("写入modinfo.lua失败: %w", err)
	}

	// 添加到数据库
	return s.addModInfoToDb(clusterName, lang, modid)
}

// GetUgcModInfo 获取UGC模组信息
func (s *ModService) GetUgcModInfo(clusterName, levelName string) ([]WorkshopItemDetail, error) {
	acfPath := s.pathResolver.GetUgcAcfPath(clusterName, levelName)
	acfWorkshops := s.parseACFFile(acfPath)

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
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok {
		return nil, errors.New("解析响应失败")
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

	return workshopItemDetails, nil
}

// DeleteUgcModFile 删除UGC模组文件
func (s *ModService) DeleteUgcModFile(clusterName, levelName, workshopId string) error {
	modFilePath := s.pathResolver.GetUgcWorkshopModPath(clusterName, levelName, workshopId)
	if fileUtils.Exists(modFilePath) {
		return fileUtils.DeleteDir(modFilePath)
	}
	return nil
}

// ===== 私有方法 =====

// parseACFFile 解析ACF文件
func (s *ModService) parseACFFile(filePath string) map[string]WorkshopItem {
	lines, err := fileUtils.ReadLnFile(filePath)
	if err != nil {
		log.Println(err)
		return nil
	}
	parsingWorkshopItemsInstalled := false
	workshopItems := make(map[string]WorkshopItem)
	var currentItemID string
	var currentItem WorkshopItem

	for _, line := range lines {
		if strings.Contains(line, "WorkshopItemsInstalled") {
			parsingWorkshopItemsInstalled = true
			continue
		}

		if strings.Contains(line, "{") && parsingWorkshopItemsInstalled {
			continue
		}

		if strings.Contains(line, "}") {
			continue
		}

		if parsingWorkshopItemsInstalled {
			replace := strings.Replace(line, "\t\t", "", -1)
			replace = strings.Replace(replace, "\"", "", -1)
			if _, err := strconv.Atoi(replace); err == nil {
				// This line contains the Workshop Item ID
				fields := strings.Fields(line)
				value := strings.Replace(fields[0], "\"", "", -1)
				currentItemID = value
			} else {
				// This line contains the Workshop Item details
				fields := strings.Fields(line)
				if len(fields) == 2 {
					key := strings.Replace(fields[0], "\"", "", -1)
					value := strings.Replace(fields[1], "\"", "", -1)
					// Remove double quotes from keys
					key = strings.ReplaceAll(key, "\"", "")
					switch key {
					case "timeupdated":
						currentItem.TimeUpdated, _ = strconv.ParseInt(value, 10, 64)
					case "manifest":
						currentItem.Manifest = strings.ReplaceAll(value, "\"", "")
					case "ugchandle":
						currentItem.Ugchandle = strings.ReplaceAll(value, "\"", "")
					}
				}
			}

			if currentItemID != "" && currentItem.TimeUpdated != 0 {
				workshopItems[currentItemID] = currentItem
				currentItemID = ""
				currentItem = WorkshopItem{}
			}
		}
	}

	return workshopItems
}

// getModInfoConfig 获取mod配置信息
func (s *ModService) getModInfoConfig(clusterName, lang, modId string) map[string]interface{} {
	// 从服务器本地读取mod信息
	if dstModInstalledPath, ok := s.getDstUcgsModsInstalledPath(clusterName, modId); ok {
		modinfoPath := filepath.Join(dstModInstalledPath, "modinfo.lua")
		if _, err := os.Stat(modinfoPath); err == nil {
			return s.readModInfo(lang, modId, modinfoPath)
		}
	}

	// 检查mod文件是否已经存在
	config, _ := s.dstConfig.GetDstConfig(clusterName)
	modDownloadPath := config.Mod_download_path
	fileUtils.CreateDirIfNotExists(modDownloadPath)

	// 下载的模组位置
	modPath := filepath.Join(modDownloadPath, "steamapps", "workshop", "content", "322330", modId)
	if _, err := os.Stat(modPath); err == nil {
		log.Println("Mod already downloaded to:", modPath)
	} else {
		// 调用 SteamCMD 命令下载 mod
		steamcmd := config.Steamcmd
		if runtime.GOOS == "windows" {
			cmd := "cd /d " + steamcmd + " && Start steamcmd.exe +login anonymous +force_install_dir " + modDownloadPath + " +workshop_download_item 322330 " + modId + " +quit"
			log.Println("正在下载模组 command:", cmd)
			_, err := shellUtils.ExecuteCommandInWin(cmd)
			if err != nil {
				log.Println("下载mod失败，请检查steamcmd路径是否配置正确", err)
				return make(map[string]interface{})
			}
		} else {
			var cmd *exec.Cmd
			if fileUtils.Exists(filepath.Join(steamcmd, "steamcmd")) {
				cmd = exec.Command(filepath.Join(steamcmd, "steamcmd"), "+login anonymous", "+force_install_dir", modDownloadPath, "+workshop_download_item 322330 "+modId, "+quit")
			} else {
				cmd = exec.Command(filepath.Join(steamcmd, "steamcmd.sh"), "+login anonymous", "+force_install_dir", modDownloadPath, "+workshop_download_item 322330 "+modId, "+quit")
			}

			log.Println("正在下载模组 command:", cmd)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Println("下载mod失败，请检查steamcmd路径是否配置正确", err)
				return make(map[string]interface{})
			}

			// 解析 SteamCMD 输出
			re := regexp.MustCompile(`Downloaded item \d+ to "([^"]+)"`)
			match := re.FindStringSubmatch(string(output))
			if len(match) < 2 {
				log.Println("Error parsing output:", string(output))
				return make(map[string]interface{})
			}
			log.Println("Mod downloaded to:", match[1])
		}
	}

	// 查找 modinfo.lua 文件
	modinfoPath := filepath.Join(modPath, "modinfo.lua")
	if _, err := os.Stat(modinfoPath); err != nil {
		log.Println("Error finding modinfo.lua:", err)
		return make(map[string]interface{})
	}
	return s.readModInfo(lang, modId, modinfoPath)
}

// getV1ModInfoConfig 从v1 mod中获取配置
func (s *ModService) getV1ModInfoConfig(clusterName, lang, modid, fileUrl string) map[string]interface{} {
	log.Println("开始下载 v1 mod，并提取 modinfo.lua 文件")
	modinfo := map[string][]byte{"modinfo": nil, "modinfo_chs": nil}
	var tmp bytes.Buffer

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", fileUrl, nil)
		if err != nil {
			log.Println(fileUrl, "下载失败", err)
			continue
		}
		client := http.Client{
			Timeout: time.Duration(10 * time.Second),
		}
		res, err := client.Do(req)
		if err != nil {
			log.Println(fileUrl, "下载失败", err)
			continue
		}
		defer res.Body.Close()
		_, err = tmp.ReadFrom(res.Body)
		if err != nil {
			log.Println(err)
			continue
		}
		break
	}

	if tmp.Len() == 0 {
		log.Println(fileUrl, "下载失败 3 次，不再尝试")
		return make(map[string]interface{})
	}

	log.Println(fileUrl, "下载成功，开始解压")
	zipReader, err := zip.NewReader(bytes.NewReader(tmp.Bytes()), int64(tmp.Len()))
	if err != nil {
		log.Println("模组zip解压失败", err)
		return make(map[string]interface{})
	}

	_ = s.unzipToDir(zipReader, filepath.Join(s.pathResolver.GetUgcModPath(clusterName), "content", "322330", modid))

	for _, file := range zipReader.File {
		switch file.Name {
		case "modinfo.lua":
			f, _ := file.Open()
			modinfoBytes, err := ioutil.ReadAll(f)
			if err != nil {
				log.Println(fileUrl, "解压 modinfo.lua 失败", err)
				continue
			}
			modinfo["modinfo"] = modinfoBytes
		case "modinfo_chs.lua":
			f, _ := file.Open()
			modinfoBytes, err := ioutil.ReadAll(f)
			if err != nil {
				log.Println(fileUrl, "解压 modinfo_chs.lua 失败", err)
				continue
			}
			modinfo["modinfo_chs"] = modinfoBytes
		}
	}

	if modinfo["modinfo"] != nil {
		return s.parseModInfoLua(lang, modid, string(modinfo["modinfo"]))
	}
	return make(map[string]interface{})
}

// getDstUcgsModsInstalledPath 获取饥荒本身modid的位置
func (s *ModService) getDstUcgsModsInstalledPath(clusterName, modid string) (string, bool) {
	config, _ := s.dstConfig.GetDstConfig(clusterName)
	var masterModFilePath, caveModFilePath string

	if config.Ugc_directory != "" {
		masterModFilePath = filepath.Join(s.pathResolver.GetUgcModPath(clusterName), "content", "322330", modid)
		caveModFilePath = filepath.Join(s.pathResolver.GetUgcModPath(clusterName), "content", "322330", modid)
	} else {
		masterModFilePath = filepath.Join(config.Force_install_dir, "ugc_mods", clusterName, "Master", "content", "322330", modid)
		caveModFilePath = filepath.Join(config.Force_install_dir, "ugc_mods", clusterName, "Caves", "content", "322330", modid)
	}

	if fileUtils.Exists(masterModFilePath) {
		return masterModFilePath, true
	}
	if fileUtils.Exists(caveModFilePath) {
		return caveModFilePath, true
	}
	return "", false
}

// readModInfo 读取modinfo.lua文件
func (s *ModService) readModInfo(lang, modId, modinfoPath string) map[string]interface{} {
	script, err := ioutil.ReadFile(modinfoPath)
	if err != nil {
		log.Println("Error reading modinfo.lua:", err)
		return make(map[string]interface{})
	}
	return s.parseModInfoLua(lang, modId, string(script))
}

// parseModInfoLua 解析modinfo.lua文件
func (s *ModService) parseModInfoLua(lang, modId, script string) map[string]interface{} {
	L := lua.NewState()
	defer L.Close()

	L.SetGlobal("locale", lua.LString(lang))
	L.SetGlobal("folder_name", lua.LString(fmt.Sprintf("workshop-%s", modId)))
	L.SetGlobal("ChooseTranslationTable", L.NewFunction(func(L *lua.LState) int {
		tbl := L.ToTable(1)
		langTbl := tbl.RawGetString(lang)
		if langTbl != lua.LNil {
			L.Push(langTbl)
		} else {
			L.Push(tbl.RawGetInt(1))
		}
		return 1
	}))

	L.DoString(script)

	global := L.Get(lua.GlobalsIndex).(*lua.LTable)
	m := make(map[string]interface{})
	global.ForEach(func(k lua.LValue, v lua.LValue) {
		if !excludeList[k.String()] && v.Type() != lua.LTFunction {
			m[k.String()] = toInterface(v)
		}
	})

	return m
}

// getVersion 从tags中获取版本号
func (s *ModService) getVersion(tags interface{}) string {
	tagList, ok := tags.([]interface{})
	if !ok {
		return ""
	}
	for _, tag := range tagList {
		tagMap, ok := tag.(map[string]interface{})
		if !ok {
			continue
		}
		tagStr, ok := tagMap["tag"].(string)
		if !ok {
			continue
		}
		if len(tagStr) > 8 && tagStr[:8] == "version:" {
			return tagStr[8:]
		}
	}
	return ""
}

// searchModInfoByWorkshopId 通过workshopId搜索mod信息
func (s *ModService) searchModInfoByWorkshopId(modID int) ModInfo {
	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	data.Set("publishedfileids[0]", strconv.Itoa(modID))
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return ModInfo{}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ModInfo{}
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return ModInfo{}
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok || len(dataList) == 0 {
		return ModInfo{}
	}

	data2 := dataList[0].(map[string]interface{})
	if data2["consumer_appid"] == nil || data2["consumer_appid"].(float64) != 322330 {
		return ModInfo{}
	}

	img := data2["preview_url"].(string)
	auth := data2["creator"].(string)
	var authorURL string
	if auth != "" {
		authorURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s/?xml=1", auth)
	}

	modId := data2["publishedfileid"].(string)
	name := data2["title"].(string)
	description := data2["file_description"].(string)
	img = fmt.Sprintf("%s?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%%23000000&letterbox=true", img)

	return ModInfo{
		ID:     modId,
		Name:   name,
		Author: authorURL,
		Desc:   description,
		Time:   int(data2["time_updated"].(float64)),
		Sub:    int(data2["subscriptions"].(float64)),
		Img:    img,
	}
}

// getLocalModInfo 获取本地mod信息
func (s *ModService) getLocalModInfo(clusterName, lang, modId string) (*model.ModInfo, error) {
	modConfigJson, _ := json.Marshal(s.getModInfoConfig(clusterName, lang, modId))
	modConfig := string(modConfigJson)

	newModInfo := &model.ModInfo{
		Auth:          "",
		ConsumerAppid: 0,
		CreatorAppid:  0,
		Description:   "",
		Modid:         modId,
		Img:           "xxx",
		LastTime:      0,
		Name:          modId,
		V:             "",
		ModConfig:     modConfig,
	}

	err := s.db.Create(newModInfo).Error
	return newModInfo, err
}

// addModInfoToDb 添加mod信息到数据库
func (s *ModService) addModInfoToDb(clusterName, lang, modid string) error {
	var modInfo *model.ModInfo
	var err error

	if !isWorkshopId(modid) {
		modInfo, err = s.getLocalModInfo(clusterName, lang, modid)
	} else {
		// 从Steam获取mod基本信息
		modInfo, err = s.getModInfo2(modid)
	}

	if err != nil {
		return fmt.Errorf("获取modinfo失败: %w", err)
	}

	// 从数据库查找是否已存在
	oldModinfo, err := s.GetModByModId(modid)
	var modConfig string
	modConfigJson, _ := json.Marshal(s.getModInfoConfig(clusterName, lang, modid))
	modConfig = string(modConfigJson)

	if err == nil && oldModinfo.Modid != "" {
		// 更新
		oldModinfo.LastTime = modInfo.LastTime
		oldModinfo.Name = modInfo.Name
		oldModinfo.Auth = modInfo.Auth
		oldModinfo.Description = modInfo.Description
		oldModinfo.Img = modInfo.Img
		oldModinfo.V = modInfo.V
		oldModinfo.ModConfig = modConfig
		return s.db.Save(oldModinfo).Error
	}

	// 新增
	modInfo.ModConfig = modConfig
	return s.db.Create(modInfo).Error
}

// getModInfo2 从Steam API获取mod基本信息
func (s *ModService) getModInfo2(modID string) (*model.ModInfo, error) {
	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	data.Set("publishedfileids[0]", modID)
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	dataList, ok := result["response"].(map[string]interface{})["publishedfiledetails"].([]interface{})
	if !ok || len(dataList) == 0 {
		return nil, errors.New("获取mod信息失败")
	}

	data2 := dataList[0].(map[string]interface{})
	img := data2["preview_url"].(string)
	auth := data2["creator"].(string)
	var authorURL string
	if auth != "" {
		authorURL = fmt.Sprintf("https://steamcommunity.com/profiles/%s/?xml=1", auth)
	}

	modId := data2["publishedfileid"].(string)
	name := data2["title"].(string)
	lastTime := data2["time_updated"].(float64)
	description := data2["file_description"].(string)
	fileUrl := data2["file_url"]
	img = fmt.Sprintf("%s?imw=64&imh=64&ima=fit&impolicy=Letterbox&imcolor=%%23000000&letterbox=true", img)
	v := s.getVersion(data2["tags"])
	creatorAppid := data2["creator_appid"].(float64)
	consumerAppid := data2["consumer_appid"].(float64)

	var fileUrlStr = ""
	if fileUrl != nil {
		fileUrlStr = fileUrl.(string)
	}

	return &model.ModInfo{
		Auth:          authorURL,
		ConsumerAppid: consumerAppid,
		CreatorAppid:  creatorAppid,
		Description:   description,
		FileUrl:       fileUrlStr,
		Modid:         modId,
		Img:           img,
		LastTime:      lastTime,
		Name:          name,
		V:             v,
	}, nil
}

// getPublishedFileDetailsBatched 批量获取mod详情
func (s *ModService) getPublishedFileDetailsBatched(workshopIds []string, batchSize int) ([]Publishedfiledetail, error) {
	var allPublishedFileDetails []Publishedfiledetail

	for i := 0; i < len(workshopIds); i += batchSize {
		end := i + batchSize
		if end > len(workshopIds) {
			end = len(workshopIds)
		}

		batch := workshopIds[i:end]
		publishedFileDetails, err := s.getPublishedFileDetailsWithGet(batch)
		if err != nil {
			return nil, err
		}

		allPublishedFileDetails = append(allPublishedFileDetails, publishedFileDetails...)
	}

	return allPublishedFileDetails, nil
}

// getPublishedFileDetailsWithGet 通过GET方式获取mod详情
func (s *ModService) getPublishedFileDetailsWithGet(workshopIds []string) ([]Publishedfiledetail, error) {
	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	for i := range workshopIds {
		data.Set("publishedfileids["+strconv.Itoa(i)+"]", workshopIds[i])
	}
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var publishedFileDetailsData struct {
		Response struct {
			Publishedfiledetails []Publishedfiledetail `json:"publishedfiledetails"`
		} `json:"response"`
	}

	err = json.NewDecoder(res.Body).Decode(&publishedFileDetailsData)
	if err != nil {
		return nil, err
	}

	return publishedFileDetailsData.Response.Publishedfiledetails, nil
}

// getPublishedFileDetails 通过POST方式获取mod详情
func (s *ModService) getPublishedFileDetails(workshopIds []string) ([]Publishedfiledetail, error) {
	url := "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("itemcount", strconv.Itoa(len(workshopIds)))
	for i := range workshopIds {
		_ = writer.WriteField("publishedfileids["+strconv.Itoa(i)+"]", workshopIds[i])
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var publishedFileDetailsData struct {
		Response struct {
			Result               int                   `json:"result"`
			Resultcount          int                   `json:"resultcount"`
			Publishedfiledetails []Publishedfiledetail `json:"publishedfiledetails"`
		} `json:"response"`
	}

	err = json.NewDecoder(res.Body).Decode(&publishedFileDetailsData)
	if err != nil {
		return nil, err
	}

	if publishedFileDetailsData.Response.Result == 1 {
		return publishedFileDetailsData.Response.Publishedfiledetails, nil
	}

	return nil, errors.New("请求失败")
}

// unzipToDir 解压zip文件到指定目录
func (s *ModService) unzipToDir(zipReader *zip.Reader, destDir string) error {
	for _, file := range zipReader.File {
		destPath := filepath.Join(destDir, file.Name)

		if !filepath.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("非法文件路径: %s", destPath)
		}

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

// ===== 辅助函数 =====

// excludeList Lua自带对象名称
var excludeList = map[string]bool{
	"_G": true, "assert": true, "collectgarbage": true, "dofile": true, "error": true,
	"getmetatable": true, "ipairs": true, "load": true, "loadfile": true, "module": true,
	"next": true, "pairs": true, "pcall": true, "print": true, "rawequal": true, "rawget": true,
	"rawset": true, "require": true, "select": true, "setmetatable": true, "tonumber": true,
	"tostring": true, "type": true, "unpack": true, "xpcall": true, "debug": true, "_VERSION": true,
	"os": true, "_GOPHER_LUA_VERSION": true, "string": true, "math": true, "io": true, "channel": true,
	"package": true, "coroutine": true, "table": true,
}

// toInterface 将Lua值转换为interface{}
func toInterface(lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(lv.(lua.LBool))
	case lua.LTNumber:
		return float64(lv.(lua.LNumber))
	case lua.LTString:
		return lv.String()
	case lua.LTTable:
		t := lv.(*lua.LTable)
		if isTableArray(t) {
			arr := make([]interface{}, t.Len())
			t.ForEach(func(i lua.LValue, v lua.LValue) {
				index := int(float64(i.(lua.LNumber)) - 1)
				if index != -1 && index < len(arr) {
					arr[index] = toInterface(v)
				}
			})
			return arr
		}
		return toMap(t)
	default:
		return lv.String()
	}
}

// toMap 将Lua table转换为map
func toMap(t *lua.LTable) map[string]interface{} {
	m := make(map[string]interface{})
	t.ForEach(func(k lua.LValue, v lua.LValue) {
		key := ""
		switch k.Type() {
		case lua.LTString:
			key = k.String()
		case lua.LTNumber:
			key = fmt.Sprintf("%g", float64(k.(lua.LNumber)))
		default:
			key = fmt.Sprintf("%v", k)
		}
		m[key] = toInterface(v)
	})
	return m
}

// isTableArray 判断Lua table是否为数组
func isTableArray(t *lua.LTable) bool {
	maxIndex := 0
	isSequential := true
	t.ForEach(func(k lua.LValue, v lua.LValue) {
		if i, ok := k.(lua.LNumber); ok {
			if i != lua.LNumber(int(i)) {
				isSequential = false
			} else if int(i) > maxIndex {
				maxIndex = int(i)
			}
		} else {
			isSequential = false
		}
	})
	return isSequential && maxIndex == t.Len()
}

// isWorkshopId 判断是否为workshop ID
func isWorkshopId(id string) bool {
	_, err := strconv.Atoi(id)
	return err == nil
}

// isModId 判断字符串是否为modID
func isModId(str string) (int, bool) {
	id, err := strconv.Atoi(str)
	return id, err == nil
}

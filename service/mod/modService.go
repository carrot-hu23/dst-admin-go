package mod

import (
	"bytes"
	"crypto/tls"
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/fileUtils"
	"encoding/json"
	"errors"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"sync"
)

type Publishedfiledetail struct {
	Publishedfileid string  `json:"publishedfileid"`
	Result          int     `json:"result"`
	Creator         string  `json:"creator"`
	CreatorAppID    int     `json:"creator_app_id"`
	ConsumerAppID   int     `json:"consumer_app_id"`
	Filename        string  `json:"filename"`
	FileURL         string  `json:"file_url"`
	HcontentFile    string  `json:"hcontent_file"`
	PreviewURL      string  `json:"preview_url"`
	HcontentPreview string  `json:"hcontent_preview"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	TimeCreated     float64 `json:"time_created"`
	TimeUpdated     float64 `json:"time_updated"`
	Visibility      int     `json:"visibility"`

	BanReason             string `json:"ban_reason"`
	Subscriptions         int    `json:"subscriptions"`
	Favorited             int    `json:"favorited"`
	LifetimeSubscriptions int    `json:"lifetime_subscriptions"`
	LifetimeFavorited     int    `json:"lifetime_favorited"`
	Views                 int    `json:"views"`
	Tags                  []struct {
		Tag string `json:"tag"`
	} `json:"tags"`
}

type PublishedFileDetailsData struct {
	Response struct {
		Result               int                   `json:"result"`
		Resultcount          int                   `json:"resultcount"`
		Publishedfiledetails []Publishedfiledetail `json:"publishedfiledetails"`
	} `json:"response"`
}

type PublishedFileDetailsDataGet struct {
	Response struct {
		Publishedfiledetails []Publishedfiledetail `json:"publishedfiledetails"`
	} `json:"response"`
}

func GetPublishedFileDetails(workshopIds []string) ([]Publishedfiledetail, error) {
	url := "https://api.steampowered.com/ISteamRemoteStorage/GetPublishedFileDetails/v1/"
	log.Println(url)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	_ = writer.WriteField("itemcount", strconv.Itoa(len(workshopIds)))

	for i := range workshopIds {
		_ = writer.WriteField("publishedfileids["+strconv.Itoa(i)+"]", workshopIds[i])
	}

	err := writer.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 创建支持TLS的http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	// 发送请求
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var publishedFileDetailsData PublishedFileDetailsData
	err = json.NewDecoder(res.Body).Decode(&publishedFileDetailsData)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if publishedFileDetailsData.Response.Result == 1 {
		return publishedFileDetailsData.Response.Publishedfiledetails, nil
	}
	return nil, errors.New("请求失败")
}

func GetPublishedFileDetailsBatched(workshopIds []string, batchSize int) ([]Publishedfiledetail, error) {
	var allPublishedFileDetails []Publishedfiledetail

	// 拆分 workshopIds 到批次
	for i := 0; i < len(workshopIds); i += batchSize {
		end := i + batchSize
		if end > len(workshopIds) {
			end = len(workshopIds)
		}

		// 获取当前批次的 workshopIds
		batch := workshopIds[i:end]

		// 调用原始函数
		publishedFileDetails, err := GetPublishedFileDetailsWithGet(batch)
		if err != nil {
			return nil, err
		}

		// 将结果添加到总结果中
		allPublishedFileDetails = append(allPublishedFileDetails, publishedFileDetails...)
	}

	return allPublishedFileDetails, nil
}

func GetPublishedFileDetailsWithGet(workshopIds []string) ([]Publishedfiledetail, error) {
	urlStr := "http://api.steampowered.com/IPublishedFileService/GetDetails/v1/"
	data := url.Values{}
	data.Set("key", steamAPIKey)
	data.Set("language", "6")
	for i := range workshopIds {
		data.Set("publishedfileids["+strconv.Itoa(i)+"]", workshopIds[i])
	}
	urlStr = urlStr + "?" + data.Encode()

	req, err := http.NewRequest("GET", urlStr, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var publishedFileDetailsData PublishedFileDetailsDataGet
	err = json.NewDecoder(res.Body).Decode(&publishedFileDetailsData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return publishedFileDetailsData.Response.Publishedfiledetails, nil
	// return nil, errors.New("请求失败")
}

// GetModModinfoLua TODO 获取模组 modinfo.lua
func GetModModinfoLua(workshopId string) string {
	return ""
}

func CheckModInfoUpdate() {
	var modInfos []model.ModInfo
	var needUpdateList []model.ModInfo
	var workshopIds []string

	db := database.DB
	db.Find(&modInfos)

	for i := range modInfos {
		workshopIds = append(workshopIds, modInfos[i].Modid)
	}

	publishedFileDetails, err := GetPublishedFileDetails(workshopIds)
	if err == nil {
		for i := range publishedFileDetails {
			publishedfiledetail := publishedFileDetails[i]
			for j := range modInfos {
				if modInfos[j].Modid == publishedfiledetail.Publishedfileid && modInfos[j].LastTime < publishedfiledetail.TimeUpdated {
					needUpdateList = append(needUpdateList, modInfos[i])
				}
			}
		}
		if len(needUpdateList) > 0 {
			for i := range needUpdateList {
				needUpdateList[i].Update = true
			}
			db.Save(&needUpdateList)
		}
	}
}

func UpdateModinfoList(lang string) {

	var modInfos []model.ModInfo
	var needUpdateList []model.ModInfo
	var workshopIds []string

	db := database.DB
	db.Find(&modInfos)

	for i := range modInfos {
		workshopIds = append(workshopIds, modInfos[i].Modid)
	}
	publishedFileDetails, err := GetPublishedFileDetailsBatched(workshopIds, 20)
	if err != nil {
		log.Panicln(err)
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
			dstConfig := dstConfigUtils.GetDstConfig()
			mod_download_path := dstConfig.Mod_download_path
			mod_path := filepath.Join(mod_download_path, "/steamapps/workshop/content/322330/", modId)
			_ = fileUtils.DeleteDir(mod_path)
			_, _, _ = SubscribeModByModId(modId, lang)
		}(i)
	}
	wg.Wait()

}

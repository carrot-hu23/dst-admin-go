package clusterUtils

import (
	"bytes"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstConfigUtils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetCluster(clusterName string) *model.Cluster {
	config := dstConfigUtils.GetDstConfig()
	cluster := model.Cluster{
		SteamCmd:        config.Steamcmd,
		ForceInstallDir: config.Force_install_dir,
		ClusterName:     config.Cluster,
		Backup:          config.Backup,
		ModDownloadPath: config.Mod_download_path,
		Bin:             config.Bin,
		Beta:            config.Beta,
		Ugc_directory:   config.Ugc_directory,
	}
	return &cluster
}

func GetClusterFromGin(ctx *gin.Context) *model.Cluster {
	config := dstConfigUtils.GetDstConfig()
	cluster := model.Cluster{
		SteamCmd:        config.Steamcmd,
		ForceInstallDir: config.Force_install_dir,
		ClusterName:     config.Cluster,
		Backup:          config.Backup,
		ModDownloadPath: config.Mod_download_path,
		Bin:             config.Bin,
		Beta:            config.Beta,
	}
	return &cluster
}

func GetDstServerInfo(clusterName string) []DstHomeInfo {

	defer func() {
		if err := recover(); err != nil {
			log.Println("查询集群房间失败:", err)
		}
	}()

	d := "{\"page\": 1,\"paginate\": 10,\"sort_type\": \"name\",\"sort_way\": 1,\"search_type\": 1,\"search_content\": \"%s\",\"mod\": 1}"
	d2 := fmt.Sprintf(d, clusterName)
	log.Println("查询: ", d2)
	data := []byte(d2)
	// 创建HTTP请求
	url := "https://dst.liuyh.com/index/serverlist/getserverlist.html"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("33333", err)
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")
	// 发送HTTP请求
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Println("2222", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 处理错误
		log.Println("1111", err)
	}
	s := string(body)
	s = s[1 : len(s)-1]
	s = strings.Replace(s, "\\", "", -1)
	fmt.Println(s)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(s), &result)
	if err != nil {
		fmt.Println(err)
	}
	if !result["success"].(bool) {
		return []DstHomeInfo{}
	}
	homeData := result["successinfo"].(map[string]interface{})["data"].([]interface{})
	if len(homeData) == 0 {
		return []DstHomeInfo{}
	}
	var homeDataList []DstHomeInfo
	for _, d := range homeData {
		row := d.([]interface{})[0].(string)
		connected := d.([]interface{})[5].(float64)
		maxConnect := d.([]interface{})[6].(float64)
		mode := d.([]interface{})[8].(string)
		mods := d.([]interface{})[9].(float64)
		name := d.([]interface{})[10].(string)
		password := d.([]interface{})[11].(float64)
		season := d.([]interface{})[14].(string)
		region := d.([]interface{})[20].(string)
		h := DstHomeInfo{
			Row:        row,
			Connected:  connected,
			MaxConnect: maxConnect,
			Mode:       mode,
			Mods:       mods,
			Name:       name,
			Password:   password,
			Season:     season,
			Region:     region,
		}
		homeDataList = append(homeDataList, h)
	}
	return homeDataList
}

type DstHomeInfo struct {
	Row        string
	Connected  float64
	MaxConnect float64
	Mode       string
	Mods       float64
	Name       string
	Password   float64
	Season     string
	Region     string
}

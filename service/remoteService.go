package service

import (
	"dst-admin-go/cache"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"dst-admin-go/vo/level"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type RemoteService struct{}

type LevelInfo struct {
	Ps                *vo.DstPsVo      `json:"Ps"`
	Status            bool             `json:"status"`
	LevelName         string           `json:"levelName"`
	IsMaster          bool             `json:"is_master"`
	Uuid              string           `json:"uuid"`
	Leveldataoverride string           `json:"leveldataoverride"`
	Modoverrides      string           `json:"modoverrides"`
	ServerIni         *level.ServerIni `json:"server_ini"`
}

type RemoteGameStatus struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data []LevelInfo `json:"data"`
}

type RemoteGameArchive struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data vo.GameArchive `json:"data"`
}

func (r *RemoteService) GetRemoteLevelStatus(cluster model.Cluster) bool {
	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d", cluster.Ip, cluster.Port)+"/api/game/8level/status", nil)
	if err != nil {
		log.Panicln("Error creating request: %v", err)
	}

	req.Header.Set("Cookie", cache.GetToken(cluster))

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应主体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var remoteGameStatus RemoteGameStatus
	err = json.Unmarshal(body, &remoteGameStatus)
	if err != nil {
		log.Panicln("Error unmarshalling JSON: %v", err)
	}

	isRun := false
	for i := range remoteGameStatus.Data {
		if remoteGameStatus.Data[i].Status {
			isRun = true
			break
		}
	}
	return isRun
}

func (r *RemoteService) GetRemoteGameArchive(cluster model.Cluster) *vo.GameArchive {

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个新的 GET 请求
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d", cluster.Ip, cluster.Port)+"/api/game/archive", nil)
	if err != nil {
		log.Panicln("Error creating request: %v", err)
	}

	req.Header.Set("Cookie", cache.GetToken(cluster))

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应主体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	var remoteGameArchive RemoteGameArchive

	err = json.Unmarshal(body, &remoteGameArchive)
	if err != nil {
		log.Panicln("Error unmarshalling JSON: %v", err)
	}

	return &remoteGameArchive.Data
}

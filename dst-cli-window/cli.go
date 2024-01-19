package dst_cli_window

import (
	"bytes"
	"dst-admin-go/config/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type dstCliClient struct {
}

var DstCliClient *dstCliClient

func (d *dstCliClient) Command(clusterName, levelName, command string) (string, error) {

	url := "http://localhost:" + global.Config.DstCliPort + "/py/dst/command"
	payload := map[string]interface{}{
		"key":     clusterName + "_" + levelName,
		"command": command,
	}

	// 将payload转换为JSON格式
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("转换为JSON失败:", err)
		return "", err
	}

	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("创建请求失败:", err)
		return "", err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("发送请求失败:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取响应失败:", err)
		return "", err
	}

	// 打印响应内容
	log.Println("响应内容:", string(body))

	return string(body), nil
}

type PyDstStatus struct {
	Status        bool    `json:"status"`
	CpuPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
}

func (d *dstCliClient) DstStatus(clusterName, levelName string) (*PyDstStatus, error) {

	url := "http://localhost:" + global.Config.DstCliPort + "/py/dst/status"
	payload := map[string]interface{}{
		"clusterName": clusterName,
		"levelName":   levelName,
	}

	// 将payload转换为JSON格式
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("转换为JSON失败:", err)
		return nil, err
	}

	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("创建请求失败:", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("发送请求失败:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取响应失败:", err)
		return nil, err
	}

	// 打印响应内容
	log.Println("响应内容:", string(body))
	var status PyDstStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		fmt.Println("解析JSON失败:", err)
		return nil, err
	}
	return &status, nil
}

func (d *dstCliClient) DstRun(key string) error {

	url := "http://localhost:" + global.Config.DstCliPort + "/py/dst/run"
	payload := map[string]interface{}{
		"key": key,
	}

	// 将payload转换为JSON格式
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("转换为JSON失败:", err)
		return err
	}

	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("创建请求失败:", err)
		return err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("发送请求失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取响应失败:", err)
		return err
	}

	// 打印响应内容
	log.Println("响应内容:", string(body))

	return nil
}

package cache

import (
	"bytes"
	"dst-admin-go/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func init() {
	TokenMemo = New(fetchToken)
}

var TokenMemo *Memo

func GenTokenKey(cluster model.Cluster) string {
	key := fmt.Sprintf("%s:%d:%s:%s", cluster.Ip, cluster.Port, cluster.Username, cluster.Password)
	return key
}

func fetchToken(key string) (interface{}, error) {
	log.Println("key", key)
	split := strings.Split(key, ":")
	if len(split) != 4 {
		return "", errors.New("key 格式错误")
	}
	ip := split[0]
	port := split[1]
	username := split[2]
	password := split[3]
	//loginURL := fmt.Sprintf("http://%s:%s/api/login", ip, port)
	//payload := map[string]interface{}{
	//	"username": username,
	//	"password": password,
	//}
	//// 将请求数据转换为 JSON 字节
	//jsonData, err := json.Marshal(payload)
	//if err != nil {
	//	fmt.Println("Failed to marshal JSON:", err)
	//	return "", err
	//}
	//// 创建请求的 Body
	//body := bytes.NewBuffer(jsonData)
	//
	//// 发送 POST 请求
	//resp, err := http.Post(loginURL, "application/json", body)
	//if err != nil {
	//	// fmt.Println("Failed to send POST request:", err)
	//	return "", err
	//}
	//token := ""
	//for _, cookie := range resp.Cookies() {
	//	if cookie.Name == "token" {
	//		token = cookie.Value
	//	}
	//}
	//log.Println("token=" + token)
	//return "token=" + token, nil

	loginURL := fmt.Sprintf("http://%s:%s/api/login", ip, port)
	payload := map[string]interface{}{
		"username": username,
		"password": password,
	}

	// 将请求数据转换为 JSON 字节
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// 创建请求的 Body
	body := bytes.NewBuffer(jsonData)

	// 创建带有超时设置的 http.Client
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 创建 POST 请求
	req, err := http.NewRequest("POST", loginURL, body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	// 获取 token
	token := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			token = cookie.Value
		}
	}
	log.Println("token=" + token)

	return "token=" + token, nil

}

func GetToken(cluster model.Cluster) string {
	value, err := TokenMemo.Get(GenTokenKey(cluster))
	if err != nil {
		log.Println(err)
		return ""
	}
	return value.(string)
}

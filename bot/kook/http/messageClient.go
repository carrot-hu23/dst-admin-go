package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type NewMessageBody struct {
	Type         int    `json:"type"`
	Content      string `json:"content"`
	TargetId     string `json:"target_id"`
	Quote        string `json:"quote"`
	Nonce        string `json:"nonce"`
	TempTargetId string `json:"temp_target_id"`
}

func CreateNewMessage(body NewMessageBody) []byte {
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}

type MessageClient struct {
	baseUrl string
	token   string
}

func NewMessageClient(baseUrl, token string) *MessageClient {
	return &MessageClient{
		baseUrl: baseUrl,
		token:   token,
	}
}

func (c *MessageClient) BaseError(message string) {

}

func (c *MessageClient) List(targetId string) (interface{}, error) {
	url := c.baseUrl + "/api/v3/message/list?target_id=" + targetId

	// 创建请求对象
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bot "+c.token)

	// 创建支持TLS的http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	// 发送请求
	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 获取响应内容
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		var data map[string]interface{}
		err = json.Unmarshal(responseBody, &data)
		return data["data"], nil
	}
	log.Println("登录失败", response.Status, response.Body)
	return nil, errors.New("获取数据失败")

}

func (c *MessageClient) Create(messageType int, targetId string, content string, quote string, nonce string, tempTargetId string) (interface{}, error) {
	url := c.baseUrl + "/api/v3/message/create"
	log.Println("url: ", url)
	// 创建请求对象
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(CreateNewMessage(NewMessageBody{
		Type:         messageType,
		TargetId:     targetId,
		Content:      content,
		Quote:        quote,
		Nonce:        nonce,
		TempTargetId: tempTargetId,
	})))
	if err != nil {
		log.Println("创建http request error", err)
		return nil, err
	}

	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bot "+c.token)

	// 创建支持TLS的http.Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 忽略证书验证
	}
	// 发送请求
	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		log.Println("http request error", err)
		return nil, err
	}
	defer response.Body.Close()

	// 获取响应内容
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("http request error", err)
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		var data map[string]interface{}
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			log.Println(err)
		}
		log.Println("data", data)
		return data, nil
	}
	log.Println("登录失败", response.Status, response.Body)
	return nil, errors.New("获取数据失败")

}

package api

import (
	"bytes"
	"dst-admin-go/vo/third"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThirdPartyApi struct {
}

// 获取饥荒的版本号
func (t *ThirdPartyApi) GetDstVersion(c *gin.Context) {

	url := "http://ver.tugos.cn/getLocalVersion"
	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		//"Content-Disposition": `attachment; filename="gopher.png"`,
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

// 获取第三方饥荒服务器
func (t *ThirdPartyApi) GetDstHomeServerList(c *gin.Context) {

	// response, err := http.Get("https://dst.liuyh.com/index/serverlist/getserverlist.html")
	// if err != nil || response.StatusCode != http.StatusOK {
	// 	c.Status(http.StatusServiceUnavailable)
	// 	return
	// }

	param := third.NewDstHomeServerParam()
	c.Bind(param)

	bytesData, err := json.Marshal(param)
	if err != nil {
		log.Println("josn 解析异常")
	}

	b_reader := bytes.NewReader(bytesData)

	url := "https://dst.liuyh.com/index/serverlist/getserverlist.html"
	req, _ := http.NewRequest("POST", url, b_reader)
	// 比如说设置个token
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		//"Content-Disposition": `attachment; filename="gopher.png"`,
		"X-Requested-With": "XMLHttpRequest",
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

// 获取第三方饥荒服务器详情
func (t *ThirdPartyApi) GetDstHomeDetailList(c *gin.Context) {

	param := third.NewDstHomeDetailParam()
	c.Bind(param)

	bytesData, err := json.Marshal(param)
	if err != nil {
		log.Println("josn 解析异常")
	}

	b_reader := bytes.NewReader(bytesData)

	url := "https://dst.liuyh.com/index/serverlist/getserverdetail.html"
	req, _ := http.NewRequest("POST", url, b_reader)
	// 比如说设置个token
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		//"Content-Disposition": `attachment; filename="gopher.png"`,
		"X-Requested-With": "XMLHttpRequest",
	}

	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}

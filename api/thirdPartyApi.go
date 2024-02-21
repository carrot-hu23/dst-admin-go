package api

import (
	"bytes"
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"dst-admin-go/vo/third"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

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
	err := c.ShouldBind(param)
	if err != nil {
		log.Println("参数解析错误", err)
	}

	query_data := map[string]any{}
	query_data["page"] = param.Page
	query_data["paginate"] = param.Paginate
	query_data["sort_type"] = param.SortType
	query_data["sort_way"] = param.SortWay
	query_data["search_type"] = param.Search_type
	if param.Search_content != "" {
		query_data["search_content"] = param.Search_content
	}
	if param.Mode != "" {
		query_data["mode"] = param.Mode
	}
	if param.Season != "" {
		query_data["season"] = param.Season
	}
	if param.Pvp != -1 {
		query_data["pvp"] = param.Pvp
	}
	if param.Mod != -1 {
		query_data["mod"] = param.Mod
	}
	if param.Password != -1 {
		query_data["password"] = param.Password
	}
	if param.World != -1 {
		query_data["world"] = param.World
	}
	if param.Playerpercent != "" {
		query_data["playerpercent"] = param.Playerpercent
	}

	bytesData, err := json.Marshal(query_data)
	log.Println("param", string(bytesData))

	if err != nil {
		log.Println("josn 解析异常")
	}

	b_reader := bytes.NewReader(bytesData)

	url := "http://dst.liuyh.com/index/serverlist/getserverlist.html"
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
	err := c.ShouldBind(param)
	if err != nil {
		log.Println("参数解析错误", err)
	}

	bytesData, err := json.Marshal(param)
	if err != nil {
		log.Println("josn 解析异常")
	}

	b_reader := bytes.NewReader(bytesData)

	url := "http://dst.liuyh.com/index/serverlist/getserverdetail.html"
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

var lobbyServer = service.LobbyServer{}

func (t *ThirdPartyApi) QueryLobbyServerDetail(ctx *gin.Context) {

	//获取查询参数
	region := ctx.Query("region")
	rowId := ctx.Query("rowId")

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: lobbyServer.QueryLobbyHomeInfo(region, rowId),
	})

}

func (t *ThirdPartyApi) GetDstHomeServerList2(ctx *gin.Context) {

	originalURL := "https://api.dstserverlist.top/api/list"
	u, err := url.Parse(originalURL)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
		return
	}

	// 构建参数
	params := url.Values{}
	params.Add("page", ctx.DefaultQuery("current", "1"))
	params.Add("pageCount", ctx.DefaultQuery("pageSize", "10"))
	params.Add("name", ctx.Query("Name"))

	// 将参数编码为查询字符串
	queryString := params.Encode()

	// 将查询字符串附加到原始URL
	u.RawQuery = queryString

	// 获取新的URL字符串
	newURL := u.String()

	req, _ := http.NewRequest("POST", newURL, nil)
	// 比如说设置个token
	req.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		//"Content-Disposition": `attachment; filename="gopher.png"`,
		//"X-Requested-With": "XMLHttpRequest",
	}
	ctx.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

}

func (t *ThirdPartyApi) GetDstHomeDetailList2(ctx *gin.Context) {

	originalURL := "https://api.dstserverlist.top/api/details/" + ctx.Query("rowId")

	req, _ := http.NewRequest("POST", originalURL, nil)
	// 比如说设置个token
	req.Header.Set("Content-Type", "application/json")

	response, err := (&http.Client{}).Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		//"Content-Disposition": `attachment; filename="gopher.png"`,
		//"X-Requested-With": "XMLHttpRequest",
	}
	ctx.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)

}

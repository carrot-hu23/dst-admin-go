package api

import (
	"bytes"
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/vo"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type ClusterApi struct{}

var clusterManager = service.ClusterManager{}

func (c *ClusterApi) GetClusterList(ctx *gin.Context) {
	clusterManager.QueryCluster(ctx)
}

func (c *ClusterApi) CreateCluster(ctx *gin.Context) {

	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", clusterModel)

	clusterManager.CreateCluster(&clusterModel)
	if clusterModel.ClusterType != "远程" {
		global.CollectMap.AddNewCollect(clusterModel.ClusterName)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) UpdateCluster(ctx *gin.Context) {
	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", clusterModel)
	log.Println("clusterModel", clusterModel)
	clusterManager.UpdateCluster(&clusterModel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) DeleteCluster(ctx *gin.Context) {

	clusterName := ctx.Query("clusterName")

	clusterModel, err := clusterManager.DeleteCluster(clusterName)
	log.Println("删除", clusterModel)
	if err != nil {
		log.Panicln("delete cluster error", err)
	}

	global.CollectMap.RemoveCollect(clusterModel.ClusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) FetchRemoteClusterList(ctx *gin.Context) {

	var payload struct {
		Ip       string `json:"ip"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}

	// 1. 登录
	loginURL := fmt.Sprintf("http://%s:%d/api/login", payload.Ip, payload.Port)
	data := map[string]interface{}{
		"username": payload.Username,
		"password": payload.Password,
	}
	// 将请求数据转换为 JSON 字节
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Panicln(err)
	}
	// 创建请求的 Body
	body := bytes.NewBuffer(jsonData)

	// 发送 POST 请求
	resp, err := http.Post(loginURL, "application/json", body)
	if err != nil {
		log.Panicln(err)
	}
	token := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "token" {
			token = cookie.Value
		}
	}
	log.Println("token=" + token)

	// 2.拿到 cookie 请求
	req2, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/cluster", payload.Ip, payload.Port), nil)
	if err != nil {
		log.Panicln(err)
	}

	// 设置 Cookie
	req2.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	// 发送请求
	resp2, err := (&http.Client{}).Do(req2)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()
	defer resp2.Body.Close()
	// 读取并解析响应 Body
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Panicln(err)
	}
	// var clusterVOList []vo.ClusterVO
	var respVo vo.Response
	if err = json.Unmarshal(body2, &respVo); err != nil {
		log.Panicln(err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: respVo.Data,
	})

}

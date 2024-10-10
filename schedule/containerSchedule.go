package schedule

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"fmt"
	"log"
	"net/http"
	"time"
)

var containerService service.ContainerService

func CollectContainerStatus() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到错误:", r)
		}
	}()

	db := database.DB
	var clusterList []model.Cluster
	db.Find(&clusterList)
	for i := range clusterList {
		cluster := clusterList[i]
		statusInfo, err := containerService.ContainerStatusInfo(cluster.ContainerId)
		if err != nil {
			cluster.Status = "error"
			log.Println("容器id", cluster.ContainerId, "获取失败")
		} else {
			log.Println("容器id", cluster.ContainerId, statusInfo.State.Status)
			cluster.Status = statusInfo.State.Status
		}
		if statusInfo.State.Status == "running" {
			status := checkContainerStatus(cluster)
			if !status {
				cluster.Status = "installing"
			}
		}
		db2 := database.DB
		db2.Save(&cluster)
	}
}

func checkContainerStatus(cluster model.Cluster) bool {
	// 要检查的URL
	url := fmt.Sprintf("%s%d", "http://127.0.0.1:", cluster.Port)
	// 设置一个1秒的超时
	timeout := 1 * time.Second
	return isServiceAvailable(url, timeout)
}
func isServiceAvailable(url string, timeout time.Duration) bool {
	client := http.Client{
		Timeout: timeout,
	}

	// 发送GET请求
	resp, err := client.Get(url)
	if err != nil {
		return false // 请求失败，服务未响应
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	return resp.StatusCode == http.StatusOK
}

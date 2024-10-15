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

var clusterManger service.ClusterManager

func CollectContainerStatus() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到错误:", r)
		}
	}()

	db := database.DB
	var clusterList []model.Cluster
	db.Where("activate = ?", true).Find(&clusterList)
	for i := range clusterList {
		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			cluster := clusterList[i]
			statusInfo, err := containerService.ContainerStatusInfo(cluster.ClusterName)
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
		}(i)

	}
}

func CheckClusterExpired() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到错误:", r)
		}
	}()

	db := database.DB
	var clusterList []model.Cluster
	db.Where("activate = ?", true).Find(&clusterList)
	for i := range clusterList {
		cluster := clusterList[i]
		if time.Now().Unix() > cluster.ExpireTime {
			cluster.Expired = true
		} else {
			cluster.Expired = false
		}
		db.Save(&cluster)

		if time.Now().Unix() > cluster.ExpireTime+3*24*60*60 {
			log.Println("正在删除卡密", cluster.Uuid)
			_, err := clusterManger.DeleteCluster(cluster.ClusterName)
			if err != nil {
				log.Println(err)
			}
		}

	}
}

func checkContainerStatus(cluster model.Cluster) bool {
	// 要检查的URL
	url := fmt.Sprintf("http://%s:%d", cluster.Ip, cluster.Port)
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

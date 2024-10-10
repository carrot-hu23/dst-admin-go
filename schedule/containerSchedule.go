package schedule

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"fmt"
	"log"
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
			status := containerService.ContainerDstInstallStatus(cluster.ContainerId)
			if !status {
				cluster.Status = "installing"
			}
		}
		db2 := database.DB
		db2.Save(&cluster)
	}
}

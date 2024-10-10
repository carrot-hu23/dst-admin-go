package schedule

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"log"
)

var containerService service.ContainerService

func CollectContainerStatus() {
	db := database.DB
	var clusterList []model.Cluster
	db.Find(&clusterList)
	for i := range clusterList {
		cluster := clusterList[i]
		statusInfo, err := containerService.ContainerStatusInfo(cluster.ContainerId)
		if err != nil {
			cluster.Status = "Error"
			log.Println("容器id", cluster.ContainerId, "获取失败")
		} else {
			log.Println("容器id", cluster.ContainerId, statusInfo.State.Status)
			cluster.Status = statusInfo.State.Status
		}
		// TODO 如果容器是运行中，并且 没有安装好，
		db2 := database.DB
		db2.Save(&cluster)
	}
}

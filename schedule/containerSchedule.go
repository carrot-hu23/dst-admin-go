package schedule

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/service"
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
		} else {
			cluster.Status = statusInfo.State.Status
		}
		// TODO 如果容器是运行中，并且 没有安装好，

	}
}

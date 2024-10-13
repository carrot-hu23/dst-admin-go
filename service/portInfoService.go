package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"fmt"
	"gorm.io/gorm"
	"sort"
)

const startPort = 20020
const endPort = 60000

type PortInfoService struct {
}

func (p *PortInfoService) GetUsedPorts(zone string) ([]int, error) {
	var usedPorts []int
	var portInfos []model.PortInfo

	db := database.DB
	result := db.Where("zone = ?", zone).Find(&portInfos)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, portInfo := range portInfos {
		usedPorts = append(usedPorts, portInfo.Port)
	}
	sort.Ints(usedPorts)
	return usedPorts, nil
}

func (p *PortInfoService) GetAvailablePorts(zone string, count int) ([]int, error) {
	usedPorts, err := p.GetUsedPorts(zone)
	if err != nil {
		return nil, err
	}
	var availablePort []int
	for port := startPort; len(availablePort) < count && port <= endPort; port++ {
		found := false
		for _, usedPort := range usedPorts {
			if usedPort == port {
				found = true
				break
			}
		}
		if !found {
			availablePort = append(availablePort, port)
		}
	}
	if len(availablePort) < count {
		return nil, fmt.Errorf("not enough available ports")
	}
	return availablePort, nil
}

func (p *PortInfoService) SaveAvailablePort(db *gorm.DB, zone, containerId string, ports []int) error {
	var portInfos []model.PortInfo
	for _, port := range ports {
		portInfos = append(portInfos, model.PortInfo{
			Zone:        zone,
			ContainerId: containerId,
			Port:        port,
		})
	}
	result := db.Create(&portInfos)
	return result.Error
}

func (p *PortInfoService) ReleasePort(db *gorm.DB, zone, containerId string) error {
	result := db.Where("zone=? AND container_id = ?", zone, containerId).Delete(&[]model.PortInfo{})
	return result.Error
}

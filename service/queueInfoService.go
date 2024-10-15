package service

import (
	"dst-admin-go/model"
	"gorm.io/gorm"
)

type QueueInfoService struct {
}

func (z *QueueInfoService) Create(db *gorm.DB, queue model.QueueInfo) error {
	result := db.Create(&queue)
	return result.Error
}

func (z *QueueInfoService) Delete(db *gorm.DB, id uint) error {
	result := db.Unscoped().Delete(&model.QueueInfo{}, id)
	return result.Error
}

func (z *QueueInfoService) FindAll(db *gorm.DB) ([]model.QueueInfo, error) {
	var queues []model.QueueInfo
	result := db.Find(&queues)
	if result.Error != nil {
		return nil, result.Error
	}
	return queues, nil
}

func (z *QueueInfoService) UpdateQueueInfo(db *gorm.DB, zoneId uint, newName string, newIp string, newPort int) error {
	var queue model.QueueInfo
	result := db.First(&queue, zoneId)
	if result.Error != nil {
		return result.Error
	}
	if newName != "" {
		queue.Name = newName
	}
	if newIp != "" {
		queue.Ip = newIp
	}
	if newPort != 0 {
		queue.Port = newPort
	}
	result = db.Save(&queue)
	return result.Error
}

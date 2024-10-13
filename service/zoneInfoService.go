package service

import (
	"dst-admin-go/model"
	"gorm.io/gorm"
)

type ZoneInfoService struct {
}

func (z *ZoneInfoService) Create(db *gorm.DB, zone model.ZoneInfo) error {
	result := db.Create(&zone)
	return result.Error
}

func (z *ZoneInfoService) Delete(db *gorm.DB, id uint) error {
	result := db.Unscoped().Delete(&model.ZoneInfo{}, id)
	return result.Error
}

func (z *ZoneInfoService) FindAll(db *gorm.DB) ([]model.ZoneInfo, error) {
	var zones []model.ZoneInfo
	result := db.Find(&zones)
	if result.Error != nil {
		return nil, result.Error
	}
	return zones, nil
}

func (z *ZoneInfoService) UpdateZone(db *gorm.DB, zoneId uint, newName string, newIp string, newPort int) error {
	var zone model.ZoneInfo
	result := db.First(&zone, zoneId)
	if result.Error != nil {
		return result.Error
	}
	if newName != "" {
		zone.Name = newName
	}
	if newIp != "" {
		zone.Ip = newIp
	}
	if newPort != 0 {
		zone.Port = newPort
	}
	result = db.Save(&zone)
	return result.Error
}

package service

import (
	"dst-admin-go/model"
	"gorm.io/gorm"
)

type LevelTemplateService struct {
}

func (l *LevelTemplateService) Create(db *gorm.DB, template model.LevelTemplate) error {
	result := db.Create(&template)
	return result.Error
}

func (l *LevelTemplateService) Delete(db *gorm.DB, id uint) error {
	result := db.Unscoped().Delete(&model.LevelTemplate{}, id)
	return result.Error
}

func (l *LevelTemplateService) UpdateZone(db *gorm.DB, zoneId uint, newName string) error {
	var zone model.ZoneInfo
	result := db.First(&zone, zoneId)
	if result.Error != nil {
		return result.Error
	}
	if newName != "" {
		zone.Name = newName
	}
	result = db.Save(&zone)
	return result.Error
}

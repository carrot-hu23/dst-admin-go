package model

import (
	"gorm.io/gorm"
)

type ZoneInfo struct {
	gorm.Model

	Name     string `json:"name"`
	ZoneCode string `gorm:"uniqueIndex" json:"zoneCode"`
}

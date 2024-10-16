package model

import (
	"gorm.io/gorm"
)

type ZoneInfo struct {
	gorm.Model

	Name     string `json:"name"`
	ZoneCode string `gorm:"uniqueIndex;size:255" json:"zoneCode"`
}

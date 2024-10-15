package model

import (
	"gorm.io/gorm"
	"time"
)

type ZoneInfo struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time

	Name     string `json:"name"`
	ZoneCode string `gorm:"uniqueIndex" json:"zoneCode"`
}

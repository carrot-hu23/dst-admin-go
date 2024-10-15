package model

import (
	"gorm.io/gorm"
	"time"
)

type QueueInfo struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string ` json:"name"`
	QueueCode string `gorm:"uniqueIndex" json:"zoneCode"`
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
}

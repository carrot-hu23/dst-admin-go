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
	QueueCode string `gorm:"uniqueIndex" json:"queueCode"`
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
}

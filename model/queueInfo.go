package model

import (
	"gorm.io/gorm"
)

type QueueInfo struct {
	gorm.Model

	Name      string ` json:"name"`
	QueueCode string `gorm:"uniqueIndex" json:"queueCode"`
	Ip        string `json:"ip"`
	Port      int    `json:"port"`
}

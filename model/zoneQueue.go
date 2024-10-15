package model

import (
	"gorm.io/gorm"
)

type ZoneQueue struct {
	gorm.Model

	ZoneCode  string `json:"zoneCode" gorm:"uniqueIndex:idx_zone_queue"`  // part of composite unique index
	QueueCode string `json:"queueCode" gorm:"uniqueIndex:idx_zone_queue"` // part of composite unique index
}

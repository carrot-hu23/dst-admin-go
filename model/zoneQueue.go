package model

import (
	"gorm.io/gorm"
)

type ZoneQueue struct {
	gorm.Model

	ZoneCode  string `json:"zoneCode" gorm:"uniqueIndex;size:255:idx_zone_queue"`  // part of composite unique index
	QueueCode string `json:"queueCode" gorm:"uniqueIndex;size:255:idx_zone_queue"` // part of composite unique index
}

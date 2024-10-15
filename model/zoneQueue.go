package model

import (
	"gorm.io/gorm"
	"time"
)

type ZoneQueue struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time

	ZoneCode  string `json:"zoneCode" gorm:"uniqueIndex:idx_zone_queue"`  // part of composite unique index
	QueueCode string `json:"queueCode" gorm:"uniqueIndex:idx_zone_queue"` // part of composite unique index
}

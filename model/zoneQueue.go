package model

import (
	"gorm.io/gorm"
)

type ZoneQueue struct {
	gorm.Model

	ZoneCode  string `gorm:"type:varchar(255);uniqueIndex:idx_zone_queue" json:"zoneCode"`  // 指定长度
	QueueCode string `gorm:"type:varchar(255);uniqueIndex:idx_zone_queue" json:"queueCode"` // 指定长度
}

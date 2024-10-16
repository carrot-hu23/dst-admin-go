package model

import (
	"gorm.io/gorm"
	"time"
)

type ZoneInfo struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time
}

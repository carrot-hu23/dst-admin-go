package model

import "gorm.io/gorm"

type BackupSnapshot struct {
	gorm.Model
	Name         string `json:"name"`
	Interval     int    `json:"interval"`
	MaxSnapshots int    `json:"maxSnapshots"`
	Enable       int    `json:"enable"`
	IsCSave      int    `json:"isCSave"`
}

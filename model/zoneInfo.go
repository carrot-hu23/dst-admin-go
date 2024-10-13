package model

import (
	"gorm.io/gorm"
	"time"
)

type ZoneInfo struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time

	Name     string `yaml:"name" json:"name"`
	ZoneCode string `yaml:"zoneCode" json:"zoneCode"`
	Ip       string `yaml:"ip" json:"ip"`
	Port     int    `yaml:"port" json:"port"`
}

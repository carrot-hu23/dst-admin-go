package model

import (
	"gorm.io/gorm"
	"time"
)

type PortInfo struct {
	gorm.Model
	CreatedAt time.Time
	UpdatedAt time.Time

	Zone        string `json:"zone"`
	Ip          string `json:"ip"`
	Port        int    `json:"port"`
	ContainerId string `json:"containerId"`
}

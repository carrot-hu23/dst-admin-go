package model

import (
	"gorm.io/gorm"
)

type PortInfo struct {
	gorm.Model

	Zone        string `json:"zone"`
	Port        int    `json:"port"`
	ContainerId string `json:"containerId"`
}

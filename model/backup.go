package model

import "gorm.io/gorm"

type Backup struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"Path"`
	Size        string `json:"size"`
	Days        string `json:"days"`
	Season      bool   `json:"season"`
}

package model

import "gorm.io/gorm"

type AutoCheck struct {
	gorm.Model
	Name         string `json:"name"`
	Enable       int    `json:"enable"`
	Announcement string `json:"announcement"`
	Times        int    `json:"times"`
	Sleep        int    `json:"sleep"`
}

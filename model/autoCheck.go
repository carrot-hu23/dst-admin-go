package model

import "gorm.io/gorm"

type AutoCheck struct {
	gorm.Model
	Name   string `json:"name"`
	Enable int    `json:"enable"`
}

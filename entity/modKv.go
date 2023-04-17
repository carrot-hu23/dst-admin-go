package entity

import "gorm.io/gorm"

type ModKv struct {
	gorm.Model
	ModId  int    `json:"modId"`
	Config string `json:"config"`
}

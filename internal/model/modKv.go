package model

import "gorm.io/gorm"

type ModKV struct {
	gorm.Model
	UserId  string `json:"userId"`
	ModId   int    `json:"modId"`
	Config  string `json:"config"`
	Version string `json:"version"`
}

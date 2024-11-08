package model

import "gorm.io/gorm"

type KV struct {
	gorm.Model
	Key   string `json:"key"`
	Value string `json:"value"`
}

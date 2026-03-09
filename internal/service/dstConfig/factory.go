package dstConfig

import (
	"gorm.io/gorm"
)

func NewDstConfig(db *gorm.DB) Config {
	dstConfig := NewOneDstConfig(db)
	return &dstConfig
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Description string `json:"description"`
	PhotoURL    string `json:"photoURL"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 逻辑删除
}

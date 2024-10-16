package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string `gorm:"uniqueIndex;size:255" json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Description string `json:"description"`
	PhotoURL    string `json:"photoURL"`

	Role string `json:"role"`
}

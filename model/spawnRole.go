package model

import "gorm.io/gorm"

type Spawn struct {
	gorm.Model
	// Id   int
	Name string
	Role string
	Time string
}

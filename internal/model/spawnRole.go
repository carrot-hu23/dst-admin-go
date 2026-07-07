package model

import "gorm.io/gorm"

type Spawn struct {
	gorm.Model
	Name        string
	Role        string
	Time        string
	ClusterName string
}

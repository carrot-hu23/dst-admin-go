package model

import "gorm.io/gorm"

type Proxy struct {
	gorm.Model
	// Id   int
	Name        string `json:"name"`
	Description string `json:"description"`
	Ip          string `json:"ip"`
	Port        string `json:"port"`
}

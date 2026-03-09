package model

import "gorm.io/gorm"

type WebLink struct {
	gorm.Model
	Title  string `json:"title"`
	Url    string `json:"url"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

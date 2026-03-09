package model

import "gorm.io/gorm"

type ModInfo struct {
	gorm.Model
	Auth          string  `json:"auth"`
	ConsumerAppid float64 `json:"consumer_appid"`
	CreatorAppid  float64 `json:"creator_appid"`
	Description   string  `json:"description"`
	FileUrl       string  `json:"file_url"`
	Modid         string  `json:"modid"`
	Img           string  `json:"img"`
	LastTime      float64 `json:"last_time"`
	ModConfig     string  `gorm:"TYPE:json" json:"mod_config"`
	Name          string  `json:"name"`
	V             string  `json:"v"`
	Update        bool    `json:"update"`
}

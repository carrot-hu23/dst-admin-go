package entity

import "gorm.io/gorm"

type ModInfo struct {
	gorm.Model
	Auth          string                 `json:"auth"`
	ConsumerAppid int64                  `json:"consumer_appid"`
	CreatorAppid  int64                  `json:"creator_appid"`
	Description   string                 `json:"description"`
	FileUrl       interface{}            `json:"file_url"`
	Modid         string                 `json:"modid"`
	Img           string                 `json:"img"`
	LastTime      int64                  `json:"last_time"`
	ModConfig     map[string]interface{} `json:"mod_config"`
	Name          string                 `json:"name"`
	V             string                 `json:"v"`
}

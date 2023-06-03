package model

import "gorm.io/gorm"

type Connect struct {
	gorm.Model
	Ip      string
	Name    string
	KuId    string
	SteamId string
	Time    string
}

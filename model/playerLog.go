package model

import "gorm.io/gorm"

type PlayerLog struct {
	gorm.Model
	Name        string `json:"name"`
	Role        string `json:"role"`
	KuId        string `json:"kuId"`
	SteamId     string `json:"steamId"`
	Time        string `json:"time"`
	Action      string `json:"action"`
	ActionDesc  string `json:"actionDesc"`
	Ip          string `json:"ip"`
	ClusterName string `json:"clusterName"`
}

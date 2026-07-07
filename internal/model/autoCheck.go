package model

import "gorm.io/gorm"

type AutoCheck struct {
	gorm.Model
	Name         string `json:"name"`
	ClusterName  string `json:"clusterName"`
	LevelName    string `json:"levelName"`
	Uuid         string `json:"uuid"`
	Enable       int    `json:"enable"`
	Announcement string `json:"announcement"`
	Times        int    `json:"times"`
	Sleep        int    `json:"sleep"`
	Interval     int    `json:"interval"`
	CheckType    string `json:"checkType"`
}

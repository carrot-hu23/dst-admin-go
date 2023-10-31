package model

import "gorm.io/gorm"

type AutoCheck struct {
	gorm.Model
	Name            string `json:"name"`
	ClusterName     string `json:"clusterName"`
	LevelName       string `json:"levelName"`
	Enable          int    `json:"enable"`
	EnableModUpdate int    `json:"enableModUpdate"`
	EnableDownCheck int    `json:"enableDownCheck"`
	Announcement    string `json:"announcement"`
	Times           int    `json:"times"`
	Sleep           int    `json:"sleep"`
	Interval        int    `json:"interval"`
}

package model

import "gorm.io/gorm"

type JobTask struct {
	gorm.Model
	ClusterName  string `json:"clusterName"`
	LevelName    string `json:"levelName"`
	Cron         string `json:"cron"`
	Category     string `json:"category"`
	Comment      string `json:"comment"`
	Announcement string `json:"announcement"`
	Sleep        int    `json:"sleep"`
	Times        int    `json:"times"`
}

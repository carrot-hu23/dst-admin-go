package model

import "gorm.io/gorm"

type JobTask struct {
	gorm.Model
	ClusterName string `json:"clusterName"`
	Cron        string `json:"cron"`
	Category    string `json:"category"`
	Comment     string `json:"comment"`
}

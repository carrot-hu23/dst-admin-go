package model

import "gorm.io/gorm"

type Regenerate struct {
	gorm.Model
	ClusterName string `json:"clusterName"`
}

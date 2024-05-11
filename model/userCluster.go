package model

import "gorm.io/gorm"

type UserCluster struct {
	gorm.Model
	UserId    int `gorm:"unique_index:idx_user_id_cluster_id" json:"userId"`
	ClusterId int `gorm:"unique_index:idx_user_id_cluster_id" json:"clusterId"`
}

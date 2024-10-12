package model

import (
	"gorm.io/gorm"
	"time"
)

type UserCluster struct {
	gorm.Model
	UserId                int  `gorm:"unique_index:idx_user_id_cluster_id" json:"userId"`
	ClusterId             int  `gorm:"unique_index:idx_user_id_cluster_id" json:"clusterId"`
	AllowAddLevel         bool `json:"allowAddLevel"`
	AllowEditingServerIni bool `json:"allowEditingServerIni"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 逻辑删除
}

package model

import (
	"gorm.io/gorm"
	"time"
)

type Cluster struct {
	gorm.Model

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // 逻辑删除

	ClusterName string `gorm:"uniqueIndex" json:"clusterName"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Uuid string `json:"uuid"`

	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`

	ContainerId string `json:"containerId"`
	Core        int    `json:"core"`
	Memory      int    `json:"memory"`
	Disk        int    `json:"disk"`
	Image       string `json:"image"`

	LevelNum   int `json:"levelNum"`
	MaxPlayers int `json:"maxPlayers"`
	MasterPort int `json:"masterPort"`
	MaxBackup  int `json:"maxBackup"`

	Status     string `json:"status"`
	ExpireTime int64  `json:"expireTime"`

	Expired bool `json:"expired"`

	Day int64 `json:"day"`

	Activate bool `json:"activate"`

	Quantity int `json:"quantity"`

	ZoneCode string `json:"zoneCode"`
	ZoneName string `json:"zoneName"`
}

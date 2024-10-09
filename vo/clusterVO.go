package vo

import "time"

type ClusterVO struct {
	Name        string `json:"name"`
	ClusterName string `gorm:"uniqueIndex" json:"clusterName"`
	Description string `json:"description"`
	Uuid        string `json:"uuid"`
	ID          uint
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Status bool `json:"status"`

	GameArchive *GameArchive `json:"gameArchive"`

	Ip              string `json:"ip"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	ClusterPassword string `json:"clusterPassword"`
}

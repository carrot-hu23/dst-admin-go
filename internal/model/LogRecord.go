package model

import "gorm.io/gorm"

type Action int

const (
	RUN Action = iota
	STOP
	NORMAL
)

type LogRecord struct {
	gorm.Model
	Action      Action `json:"action"`
	ClusterName string `json:"clusterName"`
	LevelName   string `json:"levelName"`
}

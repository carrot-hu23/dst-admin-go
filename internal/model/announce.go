package model

import "gorm.io/gorm"

type Announce struct {
	gorm.Model
	Enable       bool   `json:"enable"`
	Frequency    int64  `json:"frequency"`
	Interval     int64  `json:"interval"`
	IntervalUnit string `json:"intervalUnit"`
	Method       string `json:"method"`
	Content      string `json:"content"`
}

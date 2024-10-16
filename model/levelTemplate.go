package model

import (
	"gorm.io/gorm"
	"time"
)

type LevelTemplate struct {
	gorm.Model

	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `json:"name"`
	IconUrl     string `json:"IconUrl"`
	Description string `json:"description"`
	LevelNum    int    `json:"levelNum"`

	Modoverrides string `json:"modoverrides"`

	Leveldataoverride1 string `json:"leveldataoverride1"`
	Leveldataoverride2 string `json:"leveldataoverride2"`
	Leveldataoverride3 string `json:"leveldataoverride3"`
	Leveldataoverride4 string `json:"leveldataoverride4"`
	Leveldataoverride5 string `json:"leveldataoverride5"`
	Leveldataoverride6 string `json:"leveldataoverride6"`
	Leveldataoverride7 string `json:"leveldataoverride7"`
}

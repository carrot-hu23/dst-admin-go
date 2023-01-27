package vo

import "time"

type BackupVo struct {
	CreateTime time.Time `json:"createTime"`
	FileName   string    `json:"fileName"`
	FileSize   int64     `json:"fileSize"`
	Time       int64     `json:"time"`
}

func NewBackupVo() *BackupVo {
	return &BackupVo{}
}

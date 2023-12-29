package service

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
)

type LogRecordService struct {
}

func (l *LogRecordService) RecordLog(clusterName, levelName string, action model.Action) {

	db := database.DB

	logRecord := model.LogRecord{}
	logRecord.ClusterName = clusterName
	logRecord.Action = action
	db.Save(&logRecord)

}

func (l *LogRecordService) GetLastLog(clusterName, levelName string) *model.LogRecord {

	db := database.DB
	logRecord := model.LogRecord{}
	db.Last(&logRecord)

	return &logRecord
}

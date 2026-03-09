package database

import (
	"dst-admin-go/internal/config"
	"dst-admin-go/internal/model"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func InitDB(config *config.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.Db), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	Db = db
	err = db.AutoMigrate(
		&model.Spawn{},
		&model.PlayerLog{},
		&model.Connect{},
		&model.Regenerate{},
		&model.ModInfo{},
		&model.Cluster{},
		&model.JobTask{},
		&model.AutoCheck{},
		&model.Announce{},
		&model.WebLink{},
		&model.BackupSnapshot{},
		&model.LogRecord{},
		&model.KV{},
	)
	if err != nil {
		log.Println("AutoMigrate error", err)
	}
	return db
}

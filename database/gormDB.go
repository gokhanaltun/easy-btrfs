package database

import (
	"easy-btrfs/models"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func GetGormSqliteDb() *gorm.DB {
	if db == nil {
		dbPath := filepath.Join("/mnt/@ebtrfs/@data/", "ebtrfs.db")

		gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic(err)
		}

		db = gormDB
		migrate(db)
	}
	return db
}

func migrate(db *gorm.DB) {
	models := []interface{}{
		&models.GeneralConfig{},
		&models.SubvolumeConfig{},
		&models.Snapshot{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			panic("failed to migrate model: " + err.Error())
		}
	}
}

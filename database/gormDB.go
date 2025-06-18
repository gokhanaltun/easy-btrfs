package database

import (
	"easy-btrfs/models"
	"fmt"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var once sync.Once

func GetGormSqliteDb() (*gorm.DB, error) {
	var dbErr error
	once.Do(func() {
		dbPath := filepath.Join("/mnt/@ebtrfs/@data/", "ebtrfs.db")

		gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			dbErr = fmt.Errorf("failed to open database at %s: %w", dbPath, err)
			return
		}

		db = gormDB
		migrationErr := migrate(db)
		if migrationErr != nil {
			dbErr = migrationErr
			return
		}
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
	models := []interface{}{
		&models.GeneralConfig{},
		&models.SubvolumeConfig{},
		&models.Snapshot{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("failed to migrate model: " + err.Error())
		}
	}
	return nil
}

package repository

import (
	"easy-btrfs/database"
	"easy-btrfs/models"

	"gorm.io/gorm"
)

type SubvolumeConfigRepository interface {
	BaseRepository
}

type subvolumeConfigRepository struct {
	*BaseRepositoryStruct
}

func NewSubvolumeConfigRepository() *subvolumeConfigRepository {
	baseRepo := NewBaseRepositoryStruct(database.GetGormSqliteDb())
	return &subvolumeConfigRepository{BaseRepositoryStruct: baseRepo}
}

func (r *subvolumeConfigRepository) FindAll() (records []models.SubvolumeConfig, tx *gorm.DB) {
	tx = r.db.Find(&records)
	return records, tx
}

func (r *subvolumeConfigRepository) FindFirstByName(name string) (record models.SubvolumeConfig, tx *gorm.DB) {
	tx = r.db.Where("name = ?", name).First(&record)
	return record, tx
}

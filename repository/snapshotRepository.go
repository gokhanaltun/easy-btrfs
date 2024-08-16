package repository

import (
	"easy-btrfs/database"
	"easy-btrfs/models"

	"gorm.io/gorm"
)

type SnapshotRepository interface {
	BaseRepository
}

type snapshotRepository struct {
	*BaseRepositoryStruct
}

func NewSnapshotRepository() *snapshotRepository {
	baseRepo := NewBaseRepositoryStruct(database.GetGormSqliteDb())
	return &snapshotRepository{BaseRepositoryStruct: baseRepo}
}

func (r *snapshotRepository) FindAll() (records []models.Snapshot, tx *gorm.DB) {
	tx = r.db.Find(&records)
	return records, tx
}

func (r *snapshotRepository) FindAllByName(name string) (records []models.Snapshot, tx *gorm.DB) {
	tx = r.db.Where("name = ?", name).Find(&records)
	return records, tx
}

func (r *snapshotRepository) FindFirstById(id int) (record models.Snapshot, tx *gorm.DB) {
	tx = r.db.Where("id = ?", id).First(&record)
	return record, tx
}

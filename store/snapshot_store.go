package store

import (
	"easy-btrfs/models"

	"gorm.io/gorm"
)

type SnapshotStore struct {
	*BaseStore[models.Snapshot]
}

func NewSnapshotStore(db *gorm.DB) *SnapshotStore {
	return &SnapshotStore{BaseStore: NewBaseStore[models.Snapshot](db)}
}

func (s *SnapshotStore) FindAllByName(name string) (records []models.Snapshot, tx *gorm.DB) {
	tx = s.db.Where("name = ?", name).Find(&records)
	return records, tx
}

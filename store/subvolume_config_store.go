package store

import (
	"easy-btrfs/models"

	"gorm.io/gorm"
)

type SubvolumeConfigStore struct {
	*BaseStore[models.SubvolumeConfig]
}

func NewSubvolumeConfigStore(db *gorm.DB) *SubvolumeConfigStore {
	return &SubvolumeConfigStore{BaseStore: NewBaseStore[models.SubvolumeConfig](db)}
}

func (s *SubvolumeConfigStore) FindFirstByName(name string) (record models.SubvolumeConfig, tx *gorm.DB) {
	tx = s.db.Where("name = ?", name).First(&record)
	return record, tx
}

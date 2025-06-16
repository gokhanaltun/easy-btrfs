package store

import (
	"easy-btrfs/models"

	"gorm.io/gorm"
)

type GeneralConfigStore struct {
	*BaseStore[models.GeneralConfig]
}

func NewGeneralConfigStore(db *gorm.DB) *GeneralConfigStore {
	return &GeneralConfigStore{BaseStore: NewBaseStore[models.GeneralConfig](db)}
}

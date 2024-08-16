package repository

import (
	"easy-btrfs/database"
)

type GeneralConfigRepository interface {
	BaseRepository
}

type generalConfigRepository struct {
	*BaseRepositoryStruct
}

func NewGeneralConfigRepository() *generalConfigRepository {
	baseRepo := NewBaseRepositoryStruct(database.GetGormSqliteDb())
	return &generalConfigRepository{BaseRepositoryStruct: baseRepo}
}

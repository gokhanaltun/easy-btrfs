package store

import (
	"gorm.io/gorm"
)

type BaseStore[T any] struct {
	db *gorm.DB
}

func NewBaseStore[T any](db *gorm.DB) *BaseStore[T] {
	return &BaseStore[T]{db: db}
}

func (s *BaseStore[T]) Save(value *T) error {
	return s.db.Save(value).Error
}

func (s *BaseStore[T]) Delete(value *T) error {
	return s.db.Delete(value).Error
}

func (s *BaseStore[T]) Count() (count int64, err error) {
	var model T
	err = s.db.Model(model).Count(&count).Error
	return count, err
}

func (s *BaseStore[T]) CountByField(fieldName string, fieldValue string, model *T) (count int64, err error) {
	err = s.db.Model(model).Where(fieldName+" = ?", fieldValue).Count(&count).Error
	return count, err
}

func (s *BaseStore[T]) DeleteByField(fieldName string, fieldValue string, model *T) *gorm.DB {
	return s.db.Model(model).Where(fieldName+" = ?", fieldValue).Delete(model)
}

func (s *BaseStore[T]) FindAll() (records []T, tx *gorm.DB) {
	tx = s.db.Find(&records)
	return records, tx
}

func (s *BaseStore[T]) FindFirstById(id int) (record T, tx *gorm.DB) {
	tx = s.db.Where("id = ?", id).First(&record)
	return record, tx
}

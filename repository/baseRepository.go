package repository

import "gorm.io/gorm"

type BaseRepository interface {
	Save(value interface{}) (tx *gorm.DB)
	Delete(value interface{}) (tx *gorm.DB)
	Count(model interface{}) (count int64, tx *gorm.DB)
	CountByField(fieldName string, fieldValue string, model interface{}) (count int64, tx *gorm.DB)
	DeleteByField(fieldName string, fieldValue string, model interface{}) (tx *gorm.DB)
}

type BaseRepositoryStruct struct {
	db *gorm.DB
}

func NewBaseRepositoryStruct(db *gorm.DB) *BaseRepositoryStruct {
	return &BaseRepositoryStruct{db: db}
}

func (r *BaseRepositoryStruct) Save(value interface{}) (tx *gorm.DB) {
	tx = r.db.Save(value)
	return tx
}

func (r *BaseRepositoryStruct) Delete(value interface{}) (tx *gorm.DB) {
	tx = r.db.Delete(value)
	return tx
}

func (r *BaseRepositoryStruct) Count(model interface{}) (count int64, tx *gorm.DB) {
	tx = r.db.Model(model).Count(&count)
	return count, tx
}

func (r *BaseRepositoryStruct) CountByField(fieldName string, fieldValue string, model interface{}) (count int64, tx *gorm.DB) {
	tx = r.db.Model(model).Where(fieldName+" = ?", fieldValue).Count(&count)
	return count, tx
}

func (r *BaseRepositoryStruct) DeleteByField(fieldName string, fieldValue string, model interface{}) (tx *gorm.DB) {
	tx = r.db.Model(model).Where(fieldName+" = ?", fieldValue).Delete(&model)
	return tx
}

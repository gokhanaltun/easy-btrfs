package models

type GeneralConfig struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;"`
	Disk string `gorm:"not null"`
}

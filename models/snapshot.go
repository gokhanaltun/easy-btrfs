package models

type Snapshot struct {
	ID            int    `gorm:"primaryKey;autoIncrement;"`
	Name          string `gorm:"not null"`
	Path          string `gorm:"not null"`
	SubvolumePath string `gorm:"not null"`
	Pre           bool   `gorm:"not null"`
	PreviousFrom  int
}

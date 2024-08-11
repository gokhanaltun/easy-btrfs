package models

type SubvolumeConfig struct {
	ID            uint       `gorm:"primaryKey;autoIncrement;"`
	Name          string     `gorm:"not null"`
	SubvolumePath string     `gorm:"not null;unique;"`
	Snapshots     []Snapshot `gorm:"foreignKey:SubvolumePath;references:SubvolumePath;constraint:onDelete:NO ACTION;"`
}

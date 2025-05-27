package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	ID        uint           `gorm:"primaryKey"`
	Front     string         `gorm:"size:255;not null"`
	Back      string         `gorm:"size:255;not null"`
	DeckID    uint           `gorm:"not null;index"`
	CreatedAt time.Time      `gorm:"DEFAULT_GENERATED;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"DEFAULT_GENERATED on update CURRENT_TIMESTAMP;type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

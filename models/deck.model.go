package models

import (
	"time"

	"gorm.io/gorm"
)

type Deck struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"size:100;not null"`
	Description string         `gorm:"size:255"`
	UserID      uint           `gorm:"not null;index"`
	CreatedAt   time.Time      `gorm:"DEFAULT_GENERATED;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"DEFAULT_GENERATED on update CURRENT_TIMESTAMP;type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

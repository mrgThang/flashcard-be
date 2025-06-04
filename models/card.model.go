package models

import (
	"time"

	"gorm.io/gorm"
)

type Card struct {
	ID               int32          `gorm:"primaryKey"`
	Front            string         `gorm:"size:255;not null"`
	Back             string         `gorm:"size:255;not null"`
	DeckID           int32          `gorm:"not null;index"`
	UserID           int32          `gorm:"not null;index"`
	CreatedAt        time.Time      `gorm:"DEFAULT_GENERATED;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time      `gorm:"DEFAULT_GENERATED on update CURRENT_TIMESTAMP;type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	EasinessFactor   float32        `gorm:"not null;default:2.5"`
	StudyTime        time.Time      `gorm:"DEFAULT_GENERATED;type:datetime;default:CURRENT_TIMESTAMP"`
	RepetitionNumber int32          `gorm:"not null;default:0"`
}

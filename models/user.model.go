package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int32          `gorm:"primaryKey"`
	Name      string         `gorm:"size:100;not null"`
	Email     string         `gorm:"size:100;uniqueIndex;not null"`
	Password  string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"DEFAULT_GENERATED;type:datetime;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"DEFAULT_GENERATED on update CURRENT_TIMESTAMP;type:datetime;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

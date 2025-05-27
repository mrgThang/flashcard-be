package services

import (
	"gorm.io/gorm"
)

type Service struct {
	Config *config.Config
	DB *gorm.DB
	UserRepository repository.UserRepository
	DeckRepository repository.DeckRepository
	CardRepository repository.CardRepoditory
}
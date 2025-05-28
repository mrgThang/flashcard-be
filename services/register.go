package services

import (
	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/config"
	"github.com/mrgThang/flashcard-be/db"
	"github.com/mrgThang/flashcard-be/repositories"
)

type Service struct {
	Config         *config.Config
	DB             *gorm.DB
	UserRepository repositories.UserRepository
	DeckRepository repositories.DeckRepository
	CardRepository repositories.CardRepository
}

func NewService() *Service {
	cfg := &config.Config{}
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	db := db.MustConnectMysql(cfg.MysqlConfig)

	return &Service{
		Config:         cfg,
		DB:             db,
		UserRepository: repositories.NewUserRepository(db),
		DeckRepository: repositories.NewDeckRepository(db),
		CardRepository: repositories.NewCardRepository(db),
	}
}

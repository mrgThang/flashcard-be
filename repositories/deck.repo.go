package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/models"
)

type DeckRepository interface {
	CreateDeck(ctx context.Context, req dto.CreateDeckRequest, db ...*gorm.DB) error
	UpdateDeck(ctx context.Context, req dto.UpdateDeckRequest, db ...*gorm.DB) error
	GetDecks(ctx context.Context, req dto.GetDecksRequest, db ...*gorm.DB) ([]*models.Deck, error)
}

type deckRepositoryImpl struct {
	*gorm.DB
}

func NewDeckRepository(db *gorm.DB) DeckRepository {
	return &deckRepositoryImpl{db}
}

func (r *deckRepositoryImpl) CreateDeck(ctx context.Context, req dto.CreateDeckRequest, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	deck := models.Deck{
		Name:        req.Name,
		Description: req.Description,
		UserID:      req.UserID,
	}
	return database.WithContext(ctx).Create(&deck).Error
}

func (r *deckRepositoryImpl) UpdateDeck(ctx context.Context, req dto.UpdateDeckRequest, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	return database.WithContext(ctx).Model(&models.Deck{}).Where("id = ?", req.ID).Updates(updates).Error
}

func (r *deckRepositoryImpl) GetDecks(ctx context.Context, req dto.GetDecksRequest, dbs ...*gorm.DB) ([]*models.Deck, error) {
	database := getDb(r.DB, dbs...)
	var decks []*models.Deck
	query := database.WithContext(ctx).Model(&models.Deck{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	offset := 0
	if req.OffSet > 0 {
		offset = req.OffSet
	}
	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}
	query = query.Offset(offset).Limit(limit)
	err := query.Find(&decks).Error
	if err != nil {
		return nil, err
	}
	return decks, nil
}

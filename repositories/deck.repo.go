package repositories

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

type DeckRepository interface {
	CreateDeck(ctx context.Context, req dto.CreateDeckRequest, db ...*gorm.DB) error
	UpdateDeck(ctx context.Context, req dto.UpdateDeckRequest, db ...*gorm.DB) error
	GetDecksWithPagination(ctx context.Context, req dto.GetDecksRequest, db ...*gorm.DB) ([]*models.DeckWithStats, int64, error)
	GetDetailDeck(ctx context.Context, id int32, dbs ...*gorm.DB) (*models.DeckWithStats, error)
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

func (r *deckRepositoryImpl) GetDecksWithPagination(ctx context.Context, req dto.GetDecksRequest, dbs ...*gorm.DB) ([]*models.DeckWithStats, int64, error) {
	database := getDb(r.DB, dbs...)
	var decks []*models.DeckWithStats

	// 1. Count total decks (for pagination)
	countQuery := database.WithContext(ctx).Model(&models.Deck{})
	if req.Name != "" {
		countQuery = countQuery.Where("decks.name LIKE ?", "%"+req.Name+"%")
	}
	if req.UserID != 0 {
		countQuery = countQuery.Where("decks.user_id = ?", req.UserID)
	}
	var totalItems int64
	err := countQuery.Count(&totalItems).Error
	if err != nil {
		return nil, 0, err
	}

	// 2. Fetch decks with stats
	query := database.WithContext(ctx).Model(&models.Deck{})
	query.Select("decks.*, COUNT(cards.id) as total_cards, SUM(CASE WHEN cards.study_time < NOW() THEN 1 ELSE 0 END) as cards_left")
	if req.Name != "" {
		query = query.Where("decks.name LIKE ?", "%"+req.Name+"%")
	}
	if req.UserID != 0 {
		query = query.Where("decks.user_id = ?", req.UserID)
	}
	query = query.Joins("left join cards on decks.id = cards.deck_id")
	query = query.Group("decks.id")
	offset := constant.DefaultOffset
	if req.Page > 0 {
		offset = (req.Page - 1) * req.PageSize
	}
	limit := constant.DefaultLimit
	if req.PageSize > 0 {
		limit = req.PageSize
	}

	query = query.Offset(offset).Limit(limit)
	err = query.Find(&decks).Error
	if err != nil {
		return nil, 0, err
	}

	return decks, totalItems, nil
}

func (r *deckRepositoryImpl) GetDetailDeck(ctx context.Context, id int32, dbs ...*gorm.DB) (*models.DeckWithStats, error) {
	database := getDb(r.DB, dbs...)
	var deck models.DeckWithStats
	query := database.WithContext(ctx).Model(&models.Deck{}).Select(`
		decks.*, 
		COUNT(cards.id) as total_cards, 
		SUM(CASE WHEN cards.study_time < NOW() THEN 1 ELSE 0 END) as cards_left
	`)
	query = query.Joins("left join cards on decks.id = cards.deck_id")
	query = query.Group("decks.id")

	err := query.Where("decks.id = ?", id).First(&deck).Error
	if err != nil {
		logger.Error(fmt.Sprintf("[GetDetailDeck] Error fetching deck with ID %d", id), zap.Error(err))
		return nil, err
	}
	return &deck, nil
}

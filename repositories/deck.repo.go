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
	GetDecksWithPagination(ctx context.Context, req dto.GetDecksRequest, db ...*gorm.DB) ([]*models.Deck, int64, error)
	GetDetailDeck(ctx context.Context, id int32, dbs ...*gorm.DB) (*models.Deck, error)
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

func (r *deckRepositoryImpl) GetDecksWithPagination(ctx context.Context, req dto.GetDecksRequest, dbs ...*gorm.DB) ([]*models.Deck, int64, error) {
	database := getDb(r.DB, dbs...)
	var decks []*models.Deck
	query := database.WithContext(ctx).Model(&models.Deck{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	offset := constant.DEFAULT_OFFSET
	if req.Page > 0 {
		offset = (req.Page - 1) * req.PageSize
	}
	limit := constant.DEFAULT_LIMIT
	if req.PageSize > 0 {
		limit = req.PageSize
	}

	var totalItems int64
	err := query.Count(&totalItems).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Offset(offset).Limit(limit)
	err = query.Find(&decks).Error
	if err != nil {
		return nil, 0, err
	}

	return decks, totalItems, nil
}

func (r *deckRepositoryImpl) GetDetailDeck(ctx context.Context, id int32, dbs ...*gorm.DB) (*models.Deck, error) {
	database := getDb(r.DB, dbs...)
	var deck models.Deck
	err := database.WithContext(ctx).Where("id = ?", id).First(&deck).Error
	if err != nil {
		logger.Error(fmt.Sprintf("[GetDetailDeck] Error fetching deck with ID %d", id), zap.Error(err))
		return nil, err
	}
	return &deck, nil
}

package repositories

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

type CardRepository interface {
	CreateCard(ctx context.Context, req dto.CreateCardRequest, db ...*gorm.DB) error
	GetCards(ctx context.Context, req dto.GetCardsRequest, db ...*gorm.DB) ([]*models.Card, int64, error)
	GetDetailCard(ctx context.Context, id int32, dbs ...*gorm.DB) (*models.Card, error)
	UpdateFullCard(cardToUpdate *models.Card, dbs ...*gorm.DB) error
}

type cardRepositoryImpl struct {
	*gorm.DB
}

func NewCardRepository(db *gorm.DB) CardRepository {
	return &cardRepositoryImpl{db}
}

func (r *cardRepositoryImpl) CreateCard(ctx context.Context, req dto.CreateCardRequest, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	card := models.Card{
		UserID: req.UserID,
		Front:  req.Front,
		Back:   req.Back,
		DeckID: req.DeckID,
	}
	return database.WithContext(ctx).Create(&card).Error
}

func (r *cardRepositoryImpl) UpdateFullCard(cardToUpdate *models.Card, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	return database.Save(cardToUpdate).Error
}

func (r *cardRepositoryImpl) GetCards(ctx context.Context, req dto.GetCardsRequest, dbs ...*gorm.DB) ([]*models.Card, int64, error) {
	database := getDb(r.DB, dbs...)
	var cards []*models.Card
	query := database.WithContext(ctx).Model(&models.Card{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	}
	if req.DeckID != 0 {
		query = query.Where("deck_id = ?", req.DeckID)
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.Front != "" {
		query = query.Where("front LIKE ?", "%"+req.Front+"%")
	}
	if req.Back != "" {
		query = query.Where("back LIKE ?", "%"+req.Back+"%")
	}
	if req.StudyTimeTo != nil {
		query = query.Where("study_time <= ?", req.StudyTimeTo)
	}

	offset := constant.DefaultOffset
	if req.Page > 0 {
		offset = (req.Page - 1) * req.PageSize
	}
	limit := constant.DefaultLimit
	if req.PageSize > 0 {
		limit = req.PageSize
	}

	var totalItems int64
	err := query.Count(&totalItems).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Offset(offset).Limit(limit)
	err = query.Find(&cards).Error
	if err != nil {
		return nil, 0, err
	}
	return cards, totalItems, nil
}

func (r *cardRepositoryImpl) GetDetailCard(ctx context.Context, id int32, dbs ...*gorm.DB) (*models.Card, error) {
	database := getDb(r.DB, dbs...)
	var card models.Card
	err := database.WithContext(ctx).Model(models.Card{}).Where("id = ?", id).First(&card).Error
	if err != nil {
		logger.Error("[GetDetailCard] got error", zap.Error(err))
		return nil, err
	}
	return &card, nil
}

package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/models"
)

type CardRepository interface {
	CreateCard(ctx context.Context, req dto.CreateCardRequest, db ...*gorm.DB) error
	UpdateCard(ctx context.Context, req dto.UpdateCardRequest, db ...*gorm.DB) error
	GetCards(ctx context.Context, req dto.GetCardsRequest, db ...*gorm.DB) ([]*models.Card, error)
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
		Front:  req.Front,
		Back:   req.Back,
		DeckID: req.DeckID,
	}
	return database.WithContext(ctx).Create(&card).Error
}

func (r *cardRepositoryImpl) UpdateCard(ctx context.Context, req dto.UpdateCardRequest, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	updates := map[string]interface{}{}
	if req.Front != "" {
		updates["front"] = req.Front
	}
	if req.Back != "" {
		updates["back"] = req.Back
	}
	return database.WithContext(ctx).Model(&models.Card{}).Where("id = ?", req.ID).Updates(updates).Error
}

func (r *cardRepositoryImpl) GetCards(ctx context.Context, req dto.GetCardsRequest, dbs ...*gorm.DB) ([]*models.Card, error) {
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
	offset := 0
	if req.OffSet > 0 {
		offset = req.OffSet
	}
	limit := 10
	if req.Limit > 0 {
		limit = req.Limit
	}
	query = query.Offset(offset).Limit(limit)
	err := query.Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

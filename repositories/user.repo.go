package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest, db ...*gorm.DB) error
	GetUser(ctx context.Context, req dto.GetUserRequest, db ...*gorm.DB) (*models.User, error)
}

type userRepositoryImpl struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, req dto.CreateUserRequest, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	return database.WithContext(ctx).Create(&user).Error
}

func (r *userRepositoryImpl) GetUser(ctx context.Context, req dto.GetUserRequest, dbs ...*gorm.DB) (*models.User, error) {
	database := getDb(r.DB, dbs...)
	var user models.User
	query := database.WithContext(ctx).Model(models.User{})
	if req.ID > 0 {
		query = query.Where("id = ?", req.ID)
	}
	if len(req.Email) > 0 {
		query = query.Where("email = ?", req.Email)
	}
	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

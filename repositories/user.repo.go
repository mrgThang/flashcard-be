package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req dto.CreateUserRequest, db ...*gorm.DB) error
	UpdateUser(ctx context.Context, req dto.UpdateUserRequest, db ...*gorm.DB) error
	GetUsers(ctx context.Context, req dto.GetUserRequest, db ...*gorm.DB) (*models.User, error)
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

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, req dto.UpdateUserRequest, dbs ...*gorm.DB) error {
	database := getDb(r.DB, dbs...)
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	return database.WithContext(ctx).Model(&models.User{}).Where("id = ?", req.ID).Updates(updates).Error
}

func (r *userRepositoryImpl) GetUsers(ctx context.Context, req dto.GetUserRequest, dbs ...*gorm.DB) (*models.User, error) {
	database := getDb(r.DB, dbs...)
	var user models.User
	err := database.WithContext(ctx).First(&user, req.ID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

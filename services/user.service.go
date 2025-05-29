package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseGetUserRequest(r.Context(), r)
	if err != nil {
		logger.Error("[GetUserHandler] Invalid request parameters", zap.Error(err))
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, err := s.UserRepository.GetUser(r.Context(), *req)
	if err != nil {
		logger.Error("[GetUserHandler] UserRepository.GetUsers got error", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	response := s.parseGetUserResponse(*user)
	WriteJSONResponse(w, http.StatusOK, response)
}

func (s *Service) parseGetUserRequest(ctx context.Context, r *http.Request) (*dto.GetUserRequest, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		logger.Error("[parseGetUserRequest] Missing 'id' query parameter")
		return nil, fmt.Errorf("missing 'id' query parameter")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to convert id '%s' to int", idStr), zap.Error(err))
		return nil, err
	}
	return &dto.GetUserRequest{ID: int32(id)}, nil
}

func (s *Service) parseGetUserResponse(user models.User) dto.GetUserResponse {
	return dto.GetUserResponse{User: dto.UserItem{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}}
}

func (s *Service) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseCreateUserRequest(r.Context(), r)
	if err != nil {
		logger.Error("[CreateUserHandler] Invalid request body", zap.Error(err))
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}
	err = s.UserRepository.CreateUser(r.Context(), *req)
	if err != nil {
		logger.Error("[CreateUserHandler] UserRepository.CreateUser got error", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSONResponse(w, http.StatusOK, any(nil))
}

func (s *Service) parseCreateUserRequest(ctx context.Context, r *http.Request) (*dto.CreateUserRequest, error) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseCreateUserRequest] Failed to decode request body", zap.Error(err))
		return nil, err
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	return &req, nil
}

func (s *Service) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseUpdateUserRequest(r.Context(), r)
	if err != nil {
		logger.Error("[UpdateUserHandler] Invalid request body", zap.Error(err))
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}
	err = s.UserRepository.UpdateUser(r.Context(), *req)
	if err != nil {
		logger.Error("[UpdateUserHandler] UserRepository.UpdateUser got error", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSONResponse(w, http.StatusOK, any(nil))
}

func (s *Service) parseUpdateUserRequest(ctx context.Context, r *http.Request) (*dto.UpdateUserRequest, error) {
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseUpdateUserRequest] Failed to decode request body", zap.Error(err))
		return nil, err
	}
	return &req, nil
}

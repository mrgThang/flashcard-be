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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.UserRepository.GetUser(r.Context(), *req)
	if err != nil {
		logger.Error("[GetUserHandler] UserRepository.GetUsers", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := s.parseGetUserResponse(*user)

	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

func (s *Service) parseGetUserResponse(user models.User) ApiResponse[dto.GetUserResponse] {
	response := ApiResponse[dto.GetUserResponse]{
		Code:    http.StatusOK,
		Message: "Success",
		Data: dto.GetUserResponse{User: dto.UserItem{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}},
	}
	return response
}

func (s *Service) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseCreateUserRequest(r.Context(), r)
	if err != nil {
		logger.Error("[CreateUserHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.UserRepository.CreateUser(r.Context(), *req)
	response := s.parseCreateUserResponse(err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseCreateUserRequest(ctx context.Context, r *http.Request) (*dto.CreateUserRequest, error) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}
	return &req, nil
}

func (s *Service) parseCreateUserResponse(err error) ApiResponse[any] {
	if err != nil {
		return ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return ApiResponse[any]{
		Code:    http.StatusCreated,
		Message: "User created successfully",
		Data:    nil,
	}
}

func (s *Service) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseUpdateUserRequest(r.Context(), r)
	if err != nil {
		logger.Error("[UpdateUserHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.UserRepository.UpdateUser(r.Context(), *req)
	response := s.parseUpdateUserResponse(err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseUpdateUserRequest(ctx context.Context, r *http.Request) (*dto.UpdateUserRequest, error) {
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}
	return &req, nil
}

func (s *Service) parseUpdateUserResponse(err error) ApiResponse[any] {
	if err != nil {
		return ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return ApiResponse[any]{
		Code:    http.StatusOK,
		Message: "User updated successfully",
		Data:    nil,
	}
}

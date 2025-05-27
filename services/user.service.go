package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/models"
)

func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseGetUserRequest(r)
	if err != nil {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	user, err := s.UserRepository.GetUsers(context.Background(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := s.parseGetUserResponse(*user)

	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseGetUserRequest(r *http.Request) (dto.GetUserRequest, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return dto.GetUserRequest{}, http.ErrMissingFile
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return dto.GetUserRequest{}, err
	}
	return dto.GetUserRequest{ID: int32(id)}, nil
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
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := s.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Service) CreateUser(ctx context.Context, req dto.CreateUserRequest) error {
	return s.UserRepository.CreateUser(ctx, req)
}

func (s *Service) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := s.UpdateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) UpdateUser(ctx context.Context, req dto.UpdateUserRequest) error {
	return s.UserRepository.UpdateUser(ctx, req)
}

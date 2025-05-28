package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

func (s *Service) GetDecksHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseGetDecksRequest(r.Context(), r)
	if err != nil {
		logger.Error("[GetDecksHandler] Invalid request parameters", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decks, totalItems, err := s.DeckRepository.GetDecksWithPagination(r.Context(), *req)
	if err != nil {
		logger.Error("[GetDecksHandler] DeckRepository.GetDecks", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := s.parseGetDecksResponse(decks, dto.Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
	})
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseGetDecksRequest(ctx context.Context, r *http.Request) (*dto.GetDecksRequest, error) {
	q := r.URL.Query()
	var req dto.GetDecksRequest

	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id")
		}
		req.ID = int32(id)
	}
	req.Name = q.Get("name")
	if userIDStr := q.Get("userId"); userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid userId")
		}
		req.UserID = int32(userID)
	}
	if pageStr := q.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid page")
		}
		req.Page = page
	} else {
		req.Page = constant.DEFAULT_PAGE
	}
	if pageSizeStr := q.Get("pageSize"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pageSize")
		}
		req.PageSize = pageSize
	} else {
		req.PageSize = constant.DEFAULT_PAGE_SIZE
	}
	return &req, nil
}

func (s *Service) parseGetDecksResponse(decks []*models.Deck, pagination dto.Pagination) ApiResponse[dto.GetDecksResponse] {
	deckItems := make([]dto.DeckItem, len(decks))
	for i, deck := range decks {
		deckItems[i] = dto.DeckItem{
			ID:     deck.ID,
			Name:   deck.Name,
			UserID: deck.UserID,
		}
	}
	return ApiResponse[dto.GetDecksResponse]{
		Code:    http.StatusOK,
		Message: "Success",
		Data: dto.GetDecksResponse{
			Pagination: pagination,
			Decks:      deckItems,
		},
	}
}

func (s *Service) CreateDeckHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseCreateDeckRequest(r.Context(), r)
	if err != nil {
		logger.Error("[CreateDeckHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.DeckRepository.CreateDeck(r.Context(), *req)
	response := s.parseCreateDeckResponse(err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseCreateDeckRequest(ctx context.Context, r *http.Request) (*dto.CreateDeckRequest, error) {
	var req dto.CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}
	return &req, nil
}

func (s *Service) parseCreateDeckResponse(err error) ApiResponse[any] {
	if err != nil {
		return ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return ApiResponse[any]{
		Code:    http.StatusCreated,
		Message: "Deck created successfully",
		Data:    nil,
	}
}

func (s *Service) UpdateDeckHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseUpdateDeckRequest(r.Context(), r)
	if err != nil {
		logger.Error("[UpdateDeckHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.DeckRepository.UpdateDeck(r.Context(), *req)
	response := s.parseUpdateDeckResponse(err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseUpdateDeckRequest(ctx context.Context, r *http.Request) (*dto.UpdateDeckRequest, error) {
	var req dto.UpdateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}
	return &req, nil
}

func (s *Service) parseUpdateDeckResponse(err error) ApiResponse[any] {
	if err != nil {
		return ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return ApiResponse[any]{
		Code:    http.StatusOK,
		Message: "Deck updated successfully",
		Data:    nil,
	}
}

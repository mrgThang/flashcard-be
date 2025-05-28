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

func (s *Service) GetCardsHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseGetCardsRequest(r.Context(), r)
	if err != nil {
		logger.Error("[GetCardsHandler] Invalid request parameters", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cards, totalItems, err := s.CardRepository.GetCards(r.Context(), *req)
	if err != nil {
		logger.Error("[GetCardsHandler] CardRepository.GetCards", zap.Error(err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := s.parseGetCardsResponse(cards, dto.Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
	})
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseGetCardsRequest(ctx context.Context, r *http.Request) (*dto.GetCardsRequest, error) {
	q := r.URL.Query()
	var req dto.GetCardsRequest

	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid id")
		}
		req.ID = int32(id)
	}
	if deckIDStr := q.Get("deckId"); deckIDStr != "" {
		deckID, err := strconv.Atoi(deckIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid deckId")
		}
		req.DeckID = int32(deckID)
	}
	if userIDStr := q.Get("userId"); userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid userId")
		}
		req.UserID = int32(userID)
	}
	req.Front = q.Get("front")
	req.Back = q.Get("back")
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

func (s *Service) parseGetCardsResponse(cards []*models.Card, pagination dto.Pagination) ApiResponse[dto.GetCardsResponse] {
	cardItems := make([]dto.CardItem, len(cards))
	for i, card := range cards {
		cardItems[i] = dto.CardItem{
			ID:     card.ID,
			Front:  card.Front,
			Back:   card.Back,
			DeckID: card.DeckID,
			UserID: card.UserID,
		}
	}
	return ApiResponse[dto.GetCardsResponse]{
		Code:    http.StatusOK,
		Message: "Success",
		Data: dto.GetCardsResponse{
			Pagination: pagination,
			Cards:      cardItems,
		},
	}
}

func (s *Service) CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseCreateCardRequest(r.Context(), r)
	if err != nil {
		logger.Error("[CreateCardHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.CardRepository.CreateCard(r.Context(), *req)
	response := s.parseCreateCardResponse(err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseCreateCardRequest(ctx context.Context, r *http.Request) (*dto.CreateCardRequest, error) {
	var req dto.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}
	return &req, nil
}

func (s *Service) parseCreateCardResponse(err error) ApiResponse[any] {
	if err != nil {
		return ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return ApiResponse[any]{
		Code:    http.StatusCreated,
		Message: "Card created successfully",
		Data:    nil,
	}
}

func (s *Service) UpdateCardHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseUpdateCardRequest(r.Context(), r)
	if err != nil {
		logger.Error("[UpdateCardHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.CardRepository.UpdateCard(r.Context(), *req)
	response := s.parseUpdateCardResponse(err)
	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Service) parseUpdateCardRequest(ctx context.Context, r *http.Request) (*dto.UpdateCardRequest, error) {
	var req dto.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request body")
	}
	return &req, nil
}

func (s *Service) parseUpdateCardResponse(err error) ApiResponse[any] {
	if err != nil {
		return ApiResponse[any]{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}
	return ApiResponse[any]{
		Code:    http.StatusOK,
		Message: "Card updated successfully",
		Data:    nil,
	}
}

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
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// TODO: need logic enrich userID from token
	userId := int32(0)

	req.UserID = userId
	cards, totalItems, err := s.CardRepository.GetCards(r.Context(), *req)
	if err != nil {
		logger.Error("[GetCardsHandler] CardRepository.GetCards", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	response := s.parseGetCardsResponse(cards, dto.Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
	})
	WriteJSONResponse(w, http.StatusOK, response)
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

func (s *Service) parseGetCardsResponse(cards []*models.Card, pagination dto.Pagination) dto.GetCardsResponse {
	cardItems := make([]dto.CardItem, len(cards))
	for i, card := range cards {
		cardItems[i] = dto.CardItem{
			ID:     card.ID,
			Front:  card.Front,
			Back:   card.Back,
			DeckID: card.DeckID,
		}
	}
	return dto.GetCardsResponse{
		Pagination: pagination,
		Cards:      cardItems,
	}
}

func (s *Service) CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseCreateCardRequest(r.Context(), r)
	if err != nil {
		logger.Error("[CreateCardHandler] Invalid request body", zap.Error(err))
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	//TODO: need logic enrich userID from token
	userId := int32(0)

	deck, err := s.DeckRepository.GetDetailDeck(r.Context(), req.DeckID)
	if err != nil {
		logger.Error("[CreateCardHandler] DeckRepository.GetDetailDeck", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if deck == nil {
		logger.Error("[CreateCardHandler] Deck not found", zap.Int32("deckId", req.DeckID))
		WriteJSONError(w, http.StatusNotFound, fmt.Errorf("deck not found"))
		return
	}
	if deck.UserID != userId {
		logger.Error("[CreateCardHandler] User does not have permission to create card in this deck", zap.Int32("deckId", req.DeckID), zap.Int32("userId", req.UserID))
		WriteJSONError(w, http.StatusForbidden, fmt.Errorf("user does not have permission to create card in this deck"))
		return
	}

	req.UserID = userId
	err = s.CardRepository.CreateCard(r.Context(), *req)
	if err != nil {
		logger.Error("[CreateCardHandler] CardRepository.CreateCard", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	WriteJSONResponse(w, http.StatusCreated, any(nil))
}

func (s *Service) parseCreateCardRequest(ctx context.Context, r *http.Request) (*dto.CreateCardRequest, error) {
	var req dto.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseCreateCardRequest] Failed to decode request", zap.Error(err))
		return nil, err
	}
	if req.DeckID == 0 {
		logger.Error("[parseCreateCardRequest] DeckID is required")
		return nil, fmt.Errorf("deckId is required")
	}
	if req.Front == "" {
		logger.Error("[parseCreateCardRequest] Front is required")
		return nil, fmt.Errorf("front is required")
	}
	if req.Back == "" {
		logger.Error("[parseCreateCardRequest] Back is required")
		return nil, fmt.Errorf("back is required")
	}
	return &req, nil
}

func (s *Service) UpdateCardHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseUpdateCardRequest(r.Context(), r)
	if err != nil {
		logger.Error("[UpdateCardHandler] Invalid request body", zap.Error(err))
		WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	// TODO: need logic enrich userID from token
	userId := int32(0)

	card, err := s.CardRepository.GetDetailCard(r.Context(), req.ID)
	if err != nil {
		logger.Error("[UpdateCardHandler] CardRepository.GetDetailCard", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		logger.Error("[UpdateCardHandler] Card not found", zap.Int32("cardId", req.ID))
		WriteJSONError(w, http.StatusNotFound, fmt.Errorf("card not found"))
		return
	}
	if card.UserID != userId {
		logger.Error("[UpdateCardHandler] User does not have permission to update this card", zap.Int32("cardId", req.ID), zap.Int32("userId", userId))
		WriteJSONError(w, http.StatusForbidden, fmt.Errorf("user does not have permission to update this card"))
		return
	}

	err = s.CardRepository.UpdateCard(r.Context(), *req)
	if err != nil {
		logger.Error("[UpdateCardHandler] CardRepository.UpdateCard", zap.Error(err))
		WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	WriteJSONResponse(w, http.StatusOK, any(nil))
}

func (s *Service) parseUpdateCardRequest(ctx context.Context, r *http.Request) (*dto.UpdateCardRequest, error) {
	var req dto.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseUpdateCardRequest] Failed to decode request", zap.Error(err))
		return nil, err
	}
	if req.ID == 0 {
		logger.Error("[parseUpdateCardRequest] ID is required")
		return nil, fmt.Errorf("id is required")
	}
	return &req, nil
}

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/helpers"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

func (s *Service) GetDecksHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseGetDecksRequest(r)
	if err != nil {
		logger.Error("[GetDecksHandler] Invalid request parameters", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[GetDecksHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	req.UserID = user.ID
	decks, totalItems, err := s.DeckRepository.GetDecksWithPagination(r.Context(), *req)
	if err != nil {
		logger.Error("[GetDecksHandler] DeckRepository.GetDecks", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	response := s.parseGetDecksResponse(decks, dto.Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
	})
	helpers.WriteJSONResponse(w, http.StatusOK, response)
}

func (s *Service) parseGetDecksRequest(r *http.Request) (*dto.GetDecksRequest, error) {
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
	if pageStr := q.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("invalid page")
		}
		req.Page = page
	} else {
		req.Page = constant.DefaultPage
	}
	if pageSizeStr := q.Get("pageSize"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pageSize")
		}
		req.PageSize = pageSize
	} else {
		req.PageSize = constant.DefaultPageSize
	}
	return &req, nil
}

func (s *Service) parseGetDecksResponse(decks []*models.Deck, pagination dto.Pagination) dto.GetDecksResponse {
	deckItems := make([]dto.DeckItem, len(decks))
	for i, deck := range decks {
		deckItems[i] = dto.DeckItem{
			ID:   deck.ID,
			Name: deck.Name,
		}
	}
	return dto.GetDecksResponse{
		Pagination: pagination,
		Decks:      deckItems,
	}
}

func (s *Service) CreateDeckHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseCreateDeckRequest(r)
	if err != nil {
		logger.Error("[CreateDeckHandler] Invalid request body", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[CreateDeckHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	req.UserID = user.ID
	err = s.DeckRepository.CreateDeck(r.Context(), *req)
	if err != nil {
		logger.Error("[CreateDeckHandler] DeckRepository.CreateDeck got error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.WriteJSONResponse(w, http.StatusCreated, any(nil))
}

func (s *Service) parseCreateDeckRequest(r *http.Request) (*dto.CreateDeckRequest, error) {
	var req dto.CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseCreateDeckRequest] Decode json from req got error", zap.Error(err))
		return nil, err
	}
	if req.Name == "" {
		logger.Error("[parseCreateDeckRequest] Name is required")
		return nil, fmt.Errorf("Name is required")
	}
	return &req, nil
}

func (s *Service) UpdateDeckHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseUpdateDeckRequest(r)
	if err != nil {
		logger.Error("[UpdateDeckHandler] Invalid request body", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[UpdateDeckHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	deck, err := s.DeckRepository.GetDetailDeck(r.Context(), req.ID)
	if err != nil {
		logger.Error("[UpdateDeckHandler] DeckRepository.GetDetailDeck got error", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if deck == nil {
		logger.Error("[UpdateDeckHandler] Deck not found", zap.Int32("deckId", req.ID))
		helpers.WriteJSONError(w, http.StatusNotFound, fmt.Errorf("deck not found"))
		return
	}
	if deck.UserID != user.ID {
		logger.Error("[UpdateDeckHandler] User does not have permission to update this deck", zap.Int32("deckId", req.ID), zap.Int32("userId", user.ID))
		helpers.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("user does not have permission to update this deck"))
		return
	}

	err = s.DeckRepository.UpdateDeck(r.Context(), *req)
	if err != nil {
		logger.Error("[UpdateDeckHandler] DeckRepository.UpdateDeck got error", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSONResponse(w, http.StatusOK, any(nil))
}

func (s *Service) parseUpdateDeckRequest(r *http.Request) (*dto.UpdateDeckRequest, error) {
	var req dto.UpdateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseUpdateDeckRequest] Decode json from req got error", zap.Error(err))
		return nil, err
	}
	if req.ID == 0 {
		logger.Error("[parseUpdateDeckRequest] ID is required")
		return nil, fmt.Errorf("ID is required")
	}
	return &req, nil
}

package services

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/helpers"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

func (s *Service) GetCardsHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseGetCardsRequest(r)
	if err != nil {
		logger.Error("[GetCardsHandler] Invalid request parameters", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[GetCardsHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	req.UserID = user.ID
	cards, totalItems, err := s.CardRepository.GetCards(r.Context(), *req)
	if err != nil {
		logger.Error("[GetCardsHandler] CardRepository.GetCards", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	response := s.parseGetCardsResponse(cards, dto.Pagination{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: totalItems,
	})
	helpers.WriteJSONResponse(w, http.StatusOK, response)
}

func (s *Service) parseGetCardsRequest(r *http.Request) (*dto.GetCardsRequest, error) {
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
	if isForStudyStr := q.Get("isForStudy"); isForStudyStr != "" {
		isForStudy, err := strconv.ParseBool(isForStudyStr)
		if err != nil {
			return nil, fmt.Errorf("invalid isForStudy")
		}
		if isForStudy {
			now := time.Now()
			req.StudyTimeTo = &now
		}
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
	req, err := s.parseCreateCardRequest(r)
	if err != nil {
		logger.Error("[CreateCardHandler] Invalid request body", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[GetCardsHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	deck, err := s.DeckRepository.GetDetailDeck(r.Context(), req.DeckID)
	if err != nil {
		logger.Error("[CreateCardHandler] DeckRepository.GetDetailDeck", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if deck == nil {
		logger.Error("[CreateCardHandler] Deck not found", zap.Int32("deckId", req.DeckID))
		helpers.WriteJSONError(w, http.StatusNotFound, fmt.Errorf("deck not found"))
		return
	}
	if deck.UserID != user.ID {
		logger.Error("[CreateCardHandler] User does not have permission to create card in this deck", zap.Int32("deckId", req.DeckID), zap.Int32("userId", req.UserID))
		helpers.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("user does not have permission to create card in this deck"))
		return
	}

	req.UserID = user.ID
	err = s.CardRepository.CreateCard(r.Context(), *req)
	if err != nil {
		logger.Error("[CreateCardHandler] CardRepository.CreateCard", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusCreated, any(nil))
}

func (s *Service) parseCreateCardRequest(r *http.Request) (*dto.CreateCardRequest, error) {
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
	req, err := s.parseUpdateCardRequest(r)
	if err != nil {
		logger.Error("[UpdateCardHandler] Invalid request body", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[UpdateCardHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	card, err := s.CardRepository.GetDetailCard(r.Context(), req.ID)
	if err != nil {
		logger.Error("[UpdateCardHandler] CardRepository.GetDetailCard", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		logger.Error("[UpdateCardHandler] Card not found", zap.Int32("cardId", req.ID))
		helpers.WriteJSONError(w, http.StatusNotFound, fmt.Errorf("card not found"))
		return
	}
	if card.UserID != user.ID {
		logger.Error("[UpdateCardHandler] User does not have permission to update this card", zap.Int32("cardId", req.ID), zap.Int32("userId", user.ID))
		helpers.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("user does not have permission to update this card"))
		return
	}

	err = s.CardRepository.UpdateFullCard(card)
	if err != nil {
		logger.Error("[UpdateCardHandler] CardRepository.UpdateFullCard", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSONResponse(w, http.StatusOK, any(nil))
}

func (s *Service) parseUpdateCardRequest(r *http.Request) (*dto.UpdateCardRequest, error) {
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

func (s *Service) StudyCardHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseStudyCardRequest(r)
	if err != nil {
		logger.Error("[StudyCardHandler] Invalid request body", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[StudyCardHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}

	card, err := s.CardRepository.GetDetailCard(r.Context(), req.CardId)
	if err != nil {
		logger.Error("[StudyCardHandler] CardRepository.GetDetailCard", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	if card == nil {
		logger.Error("[StudyCardHandler] Card not found", zap.Int32("cardId", req.CardId))
		helpers.WriteJSONError(w, http.StatusNotFound, fmt.Errorf("card not found"))
		return
	}
	if card.UserID != user.ID {
		logger.Error("[StudyCardHandler] User does not have permission to study this card", zap.Int32("cardId", req.CardId), zap.Int32("userId", user.ID))
		helpers.WriteJSONError(w, http.StatusForbidden, fmt.Errorf("user does not have permission to study this card"))
		return
	}

	s.updateCardWithSm2Algorithm(card, req.QualityOfResponse)

	err = s.CardRepository.UpdateFullCard(card)
	if err != nil {
		logger.Error("[StudyCardHandler] CardRepository.UpdateFullCard got error", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	response := map[string]interface{}{
		"message": "Card studied successfully",
	}
	helpers.WriteJSONResponse(w, http.StatusOK, response)
}

func (s *Service) parseStudyCardRequest(r *http.Request) (*dto.StudyCardRequest, error) {
	var req dto.StudyCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseStudyCardRequest] Failed to decode request", zap.Error(err))
		return nil, err
	}
	if req.CardId == 0 {
		logger.Error("[parseStudyCardRequest] CardId is required")
		return nil, fmt.Errorf("cardId is required")
	}
	if req.QualityOfResponse < 0 || req.QualityOfResponse > 5 {
		logger.Error("[parseStudyCardRequest] QualityOfResponse must be between 0 and 5")
		return nil, fmt.Errorf("qualityOfResponse must be between 0 and 5")
	}
	return &req, nil
}

func (s *Service) updateCardWithSm2Algorithm(card *models.Card, qualityOfResponse int32) {
	q := qualityOfResponse
	ef := card.EasinessFactor
	n := card.RepetitionNumber
	var i int32

	ef = ef + (0.1 - (5.0-float32(q))*(0.08+(5.0-float32(q))*0.02))
	if ef < 1.3 {
		ef = 1.3
	}

	if q < 3 {
		n = 0
		i = 1
	} else {
		n += 1
		if n == 1 {
			i = 1
		} else if n == 2 {
			i = 6
		} else {
			i = int32(math.Round(float64(i) * float64(ef)))
			if i < 1 {
				i = 1
			}
		}
	}

	card.EasinessFactor = ef
	card.StudyTime = time.Now().Add(24 * time.Duration(i) * time.Hour)
	card.RepetitionNumber = n
}

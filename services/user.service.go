package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/helpers"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/models"
)

func (s *Service) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(constant.UserContextKey).(models.User)
	if !ok {
		logger.Error("[GetUserHandler] Can not get user from context")
		helpers.WriteJSONError(w, http.StatusInternalServerError, fmt.Errorf("can not get user from context"))
		return
	}
	response := s.parseGetUserResponse(user)
	helpers.WriteJSONResponse(w, http.StatusOK, response)
}

func (s *Service) parseGetUserResponse(user models.User) dto.GetUserResponse {
	return dto.GetUserResponse{User: dto.UserItem{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}}
}

func (s *Service) SignupHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseSignupRequest(r)
	if err != nil {
		logger.Error("[SignUpHandler] Invalid request body", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		logger.Error("[SignUpHandler] Hashing password gor error", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	req.Password = string(hash)

	err = s.UserRepository.CreateUser(r.Context(), *req)
	if err != nil {
		logger.Error("[SignUpHandler] UserRepository.CreateUser got error", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}
	helpers.WriteJSONResponse(w, http.StatusOK, any(nil))
}

func (s *Service) parseSignupRequest(r *http.Request) (*dto.CreateUserRequest, error) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseSignupRequest] Failed to decode request body", zap.Error(err))
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

func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := s.parseLoginRequest(r)
	if err != nil {
		logger.Error("[LoginHandler] Failed to parse request", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	getUserReq := dto.GetUserRequest{
		Email: req.Email,
	}
	user, err := s.UserRepository.GetUser(r.Context(), getUserReq)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[LoginHandler] UserRepository.GetUser not found record", zap.Error(err))
			helpers.WriteJSONError(w, http.StatusBadRequest, fmt.Errorf("email is not signed up yet"))
			return
		}
		logger.Error("[LoginHandler] UserRepository.GetUser got error", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.Error("[LoginHandler] Invalid password", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(user.ID), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})
	accessTokenString, err := accessToken.SignedString([]byte(s.Config.AccessKeySecret))
	if err != nil {
		logger.Error("[LoginHandler] Failed to create accessToken", zap.Error(err))
		helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	helpers.WriteJSONResponse(w, http.StatusOK, dto.LoginResponse{
		AccessToken: accessTokenString,
	})
}

func (s *Service) parseLoginRequest(r *http.Request) (*dto.LoginRequest, error) {
	var req *dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("[parseLoginRequest] Failed to decode request body", zap.Error(err))
		return nil, err
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	return req, nil
}

package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"github.com/mrgThang/flashcard-be/constant"
	"github.com/mrgThang/flashcard-be/dto"
	"github.com/mrgThang/flashcard-be/helpers"
	"github.com/mrgThang/flashcard-be/logger"
	"github.com/mrgThang/flashcard-be/services"
)

func AuthMiddleware(s *services.Service, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			helpers.WriteJSONError(w, http.StatusUnauthorized, fmt.Errorf("authorization header format must be Bearer {token}"))
			return
		}
		tokenString = parts[1]

		claims := &jwt.RegisteredClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
			}

			return []byte(s.Config.AccessKeySecret), nil
		})
		if err != nil {
			logger.Error("[AuthMiddleware] Failed to parse token", zap.Error(err))
			helpers.WriteJSONError(w, http.StatusUnauthorized, err)
			return
		}

		if time.Now().After(claims.ExpiresAt.Time) {
			helpers.WriteJSONResponse(w, http.StatusUnauthorized, fmt.Errorf("token expired"))
			return
		}

		userId, err := strconv.ParseInt(claims.Subject, 10, 32)
		if err != nil {
			logger.Error("[AuthMiddleware] Failed to get user id from sub", zap.Error(err))
			helpers.WriteJSONError(w, http.StatusUnauthorized, err)
			return
		}

		user, err := s.UserRepository.GetUser(r.Context(), dto.GetUserRequest{ID: int32(userId)})
		if err != nil {
			logger.Error("[AuthMiddleware] Failed to get user from sub", zap.Error(err))
			helpers.WriteJSONError(w, http.StatusInternalServerError, err)
		}

		ctx := context.WithValue(r.Context(), constant.UserContextKey, *user)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

package middlewares

import (
	"context"
	"net/http"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
)

func NewAuthMiddleware(jwtService services.JWTTokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")

			if authToken == "" {
				log.Zap.Warn("Authorization token is empty")
				http.Error(w, "User unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := jwtService.GetAndValidateToken(authToken)
			if err != nil {
				log.Zap.Warnf("failed validate Authorization token: %w", err)
				http.Error(w, "User unauthorized", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), api.KeyContext("UserID"), claims.UserID))
			next.ServeHTTP(w, r)
		})
	}
}

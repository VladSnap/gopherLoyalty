package middlewares

import (
	"context"
	"net/http"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api/auth"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")

		if authToken == "" {
			log.Zap.Warn("Authorization token is empty")
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		if verify, err := auth.VerifySignCookie(authToken); err != nil || !verify {
			log.Zap.Warnf("failed verifySignCookie: %w", err)
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		authData, err := auth.DecodeCookie(authToken)

		if err != nil {
			log.Zap.Warnf("failed decodeCookie: %w", err)
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), api.KeyContext("UserID"), authData.UserID))
		next.ServeHTTP(w, r)
	})
}

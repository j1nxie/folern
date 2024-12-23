package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/j1nxie/folern/models"
	"github.com/j1nxie/folern/utils"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, http.StatusUnauthorized, models.FolernError{Message: "missing authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.Error(w, http.StatusUnauthorized, models.FolernError{Message: "invalid authorization header"})
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, models.FolernError{Message: "invalid token"})
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

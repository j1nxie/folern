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
		var token string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 {
				token = bearerToken[1]
			}
		}

		if token == "" {
			cookie, err := r.Cookie("auth_token")
			if err == nil {
				token = cookie.Value
			}
		}

		if token == "" {
			utils.Error(w, http.StatusUnauthorized, models.FolernError{Message: "unauthorized"})
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, models.FolernError{Message: "invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middleware

import (
	"blog-app/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Error(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func GetUserID(r *http.Request) (uint, bool) {
	id, ok := r.Context().Value(UserIDKey).(uint)
	return id, ok
}

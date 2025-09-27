package handler

import (
	"net/http"

	"golang.org/x/net/context"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userID int
		var token string
		var err error

		token, err = h.s.Auth.GetTokenFromCookie(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err = h.s.Auth.ParseUserID(token)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

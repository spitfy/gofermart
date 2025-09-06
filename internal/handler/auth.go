package handler

import (
	"golang.org/x/net/context"
	"net/http"
)

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

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

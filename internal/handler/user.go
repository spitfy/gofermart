package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/middleware/auth"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if ok := validateContentType(w, r); !ok {
		return
	}
	var u user.User
	if decodeJSONBody(w, r, &u) {
		return
	}

	id, err := h.s.UserService.RegisterUser(r.Context(), u)
	if errors.Is(err, user.ErrExistsUser) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = h.s.Auth.CreateToken(w, id); err != nil {
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if ok := validateContentType(w, r); !ok {
		return
	}
	var u user.User
	if decodeJSONBody(w, r, &u) {
		return
	}
	id, err := h.s.UserService.LoginUser(r.Context(), u)
	if errors.Is(err, auth.ErrUnAuth) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = h.s.Auth.CreateToken(w, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := h.s.UserService.Balance(r.Context(), userID)
	if err != nil {
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

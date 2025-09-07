package handler

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/spitfy/gofermart/internal/helper"
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
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err = h.s.Auth.CreateToken(w, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = h.s.Auth.CreateToken(w, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	orderNum := string(body)
	if !helper.IsValidLuhn(orderNum) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.s.OrderService.AddOrder(r.Context(), userID, orderNum)
	switch {
	case errors.Is(err, order.ErrOrderAnotherUser):
		w.WriteHeader(http.StatusConflict)
		return
	case errors.Is(err, order.ErrExistsOrder):
		w.WriteHeader(http.StatusOK)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	orders, err := h.s.OrderService.ListOrders(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := h.s.UserService.Balance(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if ok := validateContentType(w, r); !ok {
		return
	}
	var wr withdraw.Request
	if decodeJSONBody(w, r, &wr) {
		return
	}
	ok, err := h.s.UserService.CanWithdraw(r.Context(), userID, wr.Sum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	err = h.s.WithdrawService.Add(r.Context(), userID, wr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ws, err := h.s.WithdrawService.List(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(ws) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(ws); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

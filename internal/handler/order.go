package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/spitfy/gofermart/internal/helper"
)

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
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
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
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
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
		log.Printf("internal server error at %s %s: %v", r.Method, r.URL.Path, err)
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}

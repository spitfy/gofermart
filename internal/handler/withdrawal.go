package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spitfy/gofermart/internal/domain/withdraw"
)

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

	err := h.s.WithdrawService.Add(r.Context(), userID, wr)
	if errors.Is(err, withdraw.ErrLowBalance) {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
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

package handler

import (
	"encoding/json"
	"errors"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/model"
	storeUser "github.com/spitfy/gofermart/internal/repository/user"
	"github.com/spitfy/gofermart/internal/service/order"
	"io"
	"mime"
	"net/http"
	"strconv"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if ok := validateContentType(w, r); !ok {
		return
	}
	var user model.User
	if decodeJSONBody(w, r, &user) {
		return
	}

	id, err := h.s.UserService.RegisterUser(r.Context(), user)
	if errors.Is(err, storeUser.ErrExistsUser) {
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
	var user model.User
	if decodeJSONBody(w, r, &user) {
		return
	}
	id, err := h.s.UserService.LoginUser(r.Context(), user)
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

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	orderNum := string(body)
	_, err = strconv.Atoi(orderNum)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	status, err := h.s.OrderService.AddOrder(r.Context(), userID, orderNum)
	switch {
	case errors.Is(err, order.ErrExistsOrderNum):
		w.WriteHeader(http.StatusConflict)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	default:
		if status == "" {
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
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
	userID, ok := r.Context().Value("userID").(int)
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
	userID, ok := r.Context().Value("userID").(int)
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

package handler

import (
	"errors"
	"github.com/spitfy/gofermart/internal/model"
	"github.com/spitfy/gofermart/internal/repository"
	"net/http"
)

/*
POST /api/user/register — регистрация пользователя;
POST /api/user/login — аутентификация пользователя;
POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
*/

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if ok := validateContentType(w, r); !ok {
		return
	}
	var user model.User
	if decodeJSONBody(w, r, &user) {
		return
	}

	id, err := h.s.UserService.RegisterUser(r.Context(), user)
	if errors.Is(err, repository.ErrExistsUser) {
		w.WriteHeader(http.StatusConflict)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err = h.s.Auth.CreateToken(w, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	if ok := validateContentType(w, r); !ok {
		return
	}
	var user model.User
	if decodeJSONBody(w, r, &user) {

		return
	}
	ok, err := h.s.UserService.LoginUser(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

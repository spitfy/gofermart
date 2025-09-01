package handler

import (
	"encoding/json"
	"errors"
	"github.com/spitfy/gofermart/internal/auth"
	"github.com/spitfy/gofermart/internal/model"
	storeUser "github.com/spitfy/gofermart/internal/repository/user"
	"github.com/spitfy/gofermart/internal/service/order"
	"io"
	"mime"
	"net/http"
	"strconv"
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
	default:
		switch status {
		case model.StatusNew:
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}
	/*
		200 — номер заказа уже был загружен этим пользователем;
		202 — новый номер заказа принят в обработку;
		400 — неверный формат запроса;
		401 — пользователь не аутентифицирован;
		409 — номер заказа уже был загружен другим пользователем;
		422 — неверный формат номера заказа;
		500 — внутренняя ошибка сервера.
	*/
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
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
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

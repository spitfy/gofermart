package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/spitfy/gofermart/internal/auth"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/service/order"
	"github.com/spitfy/gofermart/internal/service/user"
	"mime"
	"net/http"
)

func newRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", h.RegisterUser)
	r.Post("/api/user/login", h.LoginUser)
	r.Post("/api/user/orders", h.authMiddleware(h.CreateOrder))
	r.Get("/api/user/balance", h.authMiddleware(h.GetUserBalance))
	r.Get("/api/user/withdrawals", h.authMiddleware(h.ListWithdrawals))
	r.Post("/api/user/balance/withdraw", h.authMiddleware(h.WithdrawBalance))

	return r
}

type Service struct {
	Auth         *auth.AuthManager
	UserService  *user.Service
	OrderService *order.Service
}

type Handler struct {
	s Service
}

func newHandler(service Service) *Handler {
	return &Handler{
		s: service,
	}
}

func Serve(cfg *config.Config, service Service) error {
	h := newHandler(service)
	router := newRouter(h)
	server := &http.Server{
		Addr:    cfg.Handler.RunAddress,
		Handler: router,
	}

	return server.ListenAndServe()
}

var allowedContent = map[string]bool{
	"application/json":   true,
	"application/x-gzip": true,
}

func validateContentType(w http.ResponseWriter, r *http.Request) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || !allowedContent[mediaType] {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // запрещать поля, которых нет в структуре
	if err := dec.Decode(&dst); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return true
	}
	return false
}

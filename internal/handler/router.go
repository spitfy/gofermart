package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/middleware/logger"
)

func newRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", h.s.Logger.LogInfo(h.RegisterUser))
	r.Post("/api/user/login", h.s.Logger.LogInfo(h.LoginUser))
	r.Post("/api/user/orders", h.s.Logger.LogInfo(h.authMiddleware(h.CreateOrder)))
	r.Get("/api/user/orders", h.s.Logger.LogInfo(h.authMiddleware(h.ListOrders)))
	r.Get("/api/user/balance", h.s.Logger.LogInfo(h.authMiddleware(h.GetUserBalance)))
	r.Get("/api/user/withdrawals", h.s.Logger.LogInfo(h.authMiddleware(h.ListWithdrawals)))
	r.Post("/api/user/balance/withdraw", h.s.Logger.LogInfo(h.authMiddleware(h.WithdrawBalance)))

	return r
}

type Service struct {
	Auth            *auth.Service
	UserService     *user.Service
	OrderService    *order.Service
	WithdrawService *withdraw.Service
	Logger          *logger.Logger
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

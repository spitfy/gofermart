package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/service"
	"net/http"
)

func newRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", h.RegisterUser)
	r.Post("/api/user/login", h.LoginUser)
	r.Post("/api/user/orders", h.CreateOrder)
	r.Get("/api/user/balance", h.GetUserBalance)
	r.Get("/api/user/withdrawals", h.ListWithdrawals)
	r.Post("/api/user/balance/withdraw", h.WithdrawBalance)

	return r
}

func newHandler(userService *service.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

func Serve(cfg *config.Config, userService *service.UserService) error {
	h := newHandler(userService)
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

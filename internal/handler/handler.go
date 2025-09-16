package handler

import (
	"github.com/spitfy/gofermart/internal/domain/order"
	"github.com/spitfy/gofermart/internal/domain/user"
	"github.com/spitfy/gofermart/internal/domain/withdraw"
	"github.com/spitfy/gofermart/internal/middleware/auth"
	"github.com/spitfy/gofermart/internal/middleware/logger"
)

type Service struct {
	Auth            auth.Servicer
	UserService     user.Servicer
	OrderService    order.Servicer
	WithdrawService withdraw.Servicer
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

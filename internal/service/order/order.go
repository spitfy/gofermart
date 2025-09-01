package order

import (
	"context"
	"errors"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/model"
	accrualServ "github.com/spitfy/gofermart/internal/service/external/accrual"
)

var (
	ErrExistsOrderNum = errors.New("order number is exists")
)

type Service struct {
	cfg *config.Config
	s   Storer
	as  *accrualServ.Service
}

type Storer interface {
	AddOrder(ctx context.Context, order model.Order) (model.Order, error)
	ListOrders(ctx context.Context, UserID int) ([]model.Order, error)
}

func NewService(cfg *config.Config, store Storer, as *accrualServ.Service) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
		as:  as,
	}
}

func (s *Service) AddOrder(ctx context.Context, userID int, orderNumber string) (model.OrderStatus, error) {
	order := model.Order{
		UserID:      userID,
		OrderNumber: orderNumber,
		Status:      model.StatusNew,
	}
	o, err := s.s.AddOrder(ctx, order)
	if err != nil {
		return "", err
	}
	if o.UserID != order.UserID {
		return "", ErrExistsOrderNum
	}
	go s.as.Call(userID, orderNumber)
	return o.Status, nil
}

func (s *Service) ListOrders(ctx context.Context, userID int) ([]model.Order, error) {
	return s.s.ListOrders(ctx, userID)
}

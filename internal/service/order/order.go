package order

import (
	"context"
	"errors"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/model"
)

var (
	ErrExistsOrderNum = errors.New("order number is exists")
)

type Service struct {
	cfg *config.Config
	s   Storer
}

type Storer interface {
	AddOrder(ctx context.Context, order model.Order) (model.Order, error)
}

func NewService(cfg *config.Config, store Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

func (s *Service) AddOrder(ctx context.Context, userID int, OrderNumber string) (model.OrderStatus, error) {
	order := model.Order{
		UserID:      userID,
		OrderNumber: OrderNumber,
		Status:      model.StatusNew,
	}
	o, err := s.s.AddOrder(ctx, order)
	if err != nil {
		return "", err
	}
	if o.UserID != order.UserID {
		return "", ErrExistsOrderNum
	}
	return o.Status, nil
}

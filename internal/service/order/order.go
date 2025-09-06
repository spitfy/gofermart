package order

import (
	"context"
	"errors"
	"github.com/spitfy/gofermart/internal/config"
	"github.com/spitfy/gofermart/internal/model"
	"github.com/spitfy/gofermart/internal/repository/order"
	accrualServ "github.com/spitfy/gofermart/internal/service/external/accrual"
	"log"
)

var (
	ErrExistsOrder      = errors.New("order number is exists")
	ErrOrderAnotherUser = errors.New("order is exists by other user")
)

type Service struct {
	cfg *config.Config
	s   order.Storer
	as  *accrualServ.Service
}

func NewService(cfg *config.Config, store order.Storer, as *accrualServ.Service) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
		as:  as,
	}
}

func (s *Service) AddOrder(ctx context.Context, userID int, number string) error {
	m := model.Order{
		UserID: userID,
		Number: number,
	}
	o, err := s.s.AddOrder(ctx, m)
	if errors.Is(err, order.ErrUniqueNum) {
		if o.UserID != m.UserID {
			return ErrOrderAnotherUser
		}
		return ErrExistsOrder
	}
	if err != nil {
		return err
	}
	log.Println("===========AddOrder=======")
	go s.as.Call(userID, number)
	return nil
}

func (s *Service) ListOrders(ctx context.Context, userID int) ([]model.Order, error) {
	return s.s.ListOrders(ctx, userID)
}

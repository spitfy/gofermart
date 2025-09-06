package order

import (
	"context"
	"errors"
	"github.com/spitfy/gofermart/internal/config"
	accrualServ "github.com/spitfy/gofermart/internal/service/external/accrual"
)

var (
	ErrExistsOrder      = errors.New("order number is exists")
	ErrOrderAnotherUser = errors.New("order is exists by other user")
)

type Service struct {
	cfg *config.Config
	s   Storer
	as  *accrualServ.Service
}

func NewService(cfg *config.Config, store Storer, as *accrualServ.Service) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
		as:  as,
	}
}

func (s *Service) AddOrder(ctx context.Context, userID int, number string) error {
	m := Order{
		UserID: userID,
		Number: number,
	}
	o, err := s.s.AddOrder(ctx, m)
	if errors.Is(err, ErrUniqueNum) {
		if o.UserID != m.UserID {
			return ErrOrderAnotherUser
		}
		return ErrExistsOrder
	}
	if err != nil {
		return err
	}
	go s.as.Call(userID, number)
	return nil
}

func (s *Service) ListOrders(ctx context.Context, userID int) ([]Order, error) {
	return s.s.ListOrders(ctx, userID)
}

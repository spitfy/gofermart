package order

import (
	"context"
	"errors"
	"github.com/spitfy/gofermart/internal/config"
	accrualServ "github.com/spitfy/gofermart/internal/service/external/accrual"
	"runtime"
)

var (
	ErrExistsOrder      = errors.New("order number is exists")
	ErrOrderAnotherUser = errors.New("order is exists by other user")
)

type Service struct {
	cfg    *config.Config
	s      Storer
	as     *accrualServ.Service
	sendCh chan orderSend
}

func NewService(cfg *config.Config, store Storer, as *accrualServ.Service) *Service {
	s := Service{
		cfg:    cfg,
		s:      store,
		as:     as,
		sendCh: make(chan orderSend),
	}

	maxProcs := runtime.GOMAXPROCS(0)
	for i := 0; i < maxProcs; i++ {
		go s.runSendWorker()
	}

	return &s
}

func (s *Service) runSendWorker() {
	for os := range s.sendCh {
		s.as.Call(os.userID, os.number)
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
	s.sendCh <- orderSend{
		userID: userID,
		number: number,
	}
	return nil
}

func (s *Service) ListOrders(ctx context.Context, userID int) ([]Order, error) {
	return s.s.ListOrders(ctx, userID)
}

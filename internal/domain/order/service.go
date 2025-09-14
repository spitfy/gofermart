package order

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"github.com/spitfy/gofermart/internal/config"
	accrualServ "github.com/spitfy/gofermart/internal/service/external/accrual"
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

type Servicer interface {
	runSendWorker(ctx context.Context, wg *sync.WaitGroup)
	AddOrder(ctx context.Context, userID int, number string) error
	ListOrders(ctx context.Context, userID int) ([]Order, error)
}

func NewService(cfg *config.Config, store Storer, as *accrualServ.Service) *Service {
	s := Service{
		cfg:    cfg,
		s:      store,
		as:     as,
		sendCh: make(chan orderSend),
	}

	maxProcs := runtime.GOMAXPROCS(0)
	as.Wg.Add(maxProcs)
	for i := 0; i < maxProcs; i++ {
		go s.runSendWorker(as.Ctx, as.Wg)
	}

	return &s
}

func (s *Service) runSendWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case os, ok := <-s.sendCh:
			if !ok {
				return
			}
			s.as.Call(os.userID, os.number)
		case <-ctx.Done():
			return
		}
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

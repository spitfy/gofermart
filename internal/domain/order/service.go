package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/cenkalti/backoff/v5"

	"go.uber.org/atomic"

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
	as     accrualServ.Servicer
	sendCh chan orderSend
	log    *zap.Logger
}

//go:generate mockgen -destination=servicer_mock.go -package=order github.com/spitfy/gofermart/internal/domain/order Servicer
type Servicer interface {
	runSendWorker(ctx context.Context, wg *sync.WaitGroup, await atomic.Time)
	AddOrder(ctx context.Context, userID int, number string) error
	ListOrders(ctx context.Context, userID int) ([]Order, error)
}

func NewService(cfg *config.Config, store Storer, as accrualServ.Servicer, logger *zap.Logger) *Service {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &Service{
		cfg:    cfg,
		s:      store,
		as:     as,
		sendCh: make(chan orderSend, 100),
		log:    logger,
	}
}

func (s *Service) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.cfg.Accrual.Interval) * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.listForAccrual(s.as.Context()); err != nil {
					log.Println("Query error:", err)
				}
			case <-quit:
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	maxProcs := runtime.GOMAXPROCS(0)
	s.as.WaitGroup().Add(maxProcs)
	for i := 0; i < maxProcs; i++ {
		go s.runSendWorker(ctx, s.as.WaitGroup(), s.as.AwaitTime())
	}
}

func (s *Service) listForAccrual(ctx context.Context) error {
	orders, err := s.s.listForAccrual(ctx)
	if err != nil {
		return err
	}
	for _, o := range orders {
		s.sendCh <- o
	}
	return nil
}

func (s *Service) runSendWorker(ctx context.Context, wg *sync.WaitGroup, await atomic.Time) {
	defer wg.Done()
	for {
		select {
		case os, ok := <-s.sendCh:
			if !ok {
				return
			}
			now := time.Now()
			t := await.Load()
			if now.Before(t) {
				timer := time.NewTimer(t.Sub(now))
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return
				}
			}
			s.as.Call(os.userID, os.number)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) AddOrder(ctx context.Context, userID int, number string) error {
	order := Order{UserID: userID, Number: number}

	result, err := backoff.Retry(ctx,
		func() (Order, error) { return s.s.addOrder(ctx, order) },
		backoff.WithMaxTries(s.cfg.DB.MaxRetries),
	)

	switch {
	case errors.Is(err, ErrUniqueNum):
		if result.UserID != userID {
			s.log.Warn("order conflict", zap.String("number", number), zap.Int("owner", result.UserID))
			return ErrOrderAnotherUser
		}
		return ErrExistsOrder

	case err != nil:
		return fmt.Errorf("order creation failed: %w", err) // Оборачиваем для контекста
	}
	return nil
}

func (s *Service) ListOrders(ctx context.Context, userID int) ([]Order, error) {
	return s.s.listOrders(ctx, userID)
}

package withdraw

import (
	"context"

	"github.com/cenkalti/backoff/v5"

	"github.com/spitfy/gofermart/internal/config"
)

type Service struct {
	cfg *config.Config
	s   Storer
}

//go:generate mockgen -destination=servicer_mock.go -package=withdraw github.com/spitfy/gofermart/internal/domain/withdraw Servicer
type Servicer interface {
	Add(ctx context.Context, userID int, req Request) error
	List(ctx context.Context, userID int) ([]Withdraw, error)
}

func NewService(cfg *config.Config, store Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

func (s *Service) Add(ctx context.Context, userID int, req Request) error {
	m := Withdraw{
		UserID: userID,
		Order:  req.Order,
		Amount: req.Sum,
	}
	add := func() (struct{}, error) {
		err := s.s.Add(ctx, m)
		return struct{}{}, err
	}
	_, err := backoff.Retry(ctx, add, backoff.WithBackOff(backoff.NewExponentialBackOff()), backoff.WithMaxTries(s.cfg.DB.MaxRetries))

	return err
}

func (s *Service) List(ctx context.Context, userID int) ([]Withdraw, error) {
	return s.s.List(ctx, userID)
}

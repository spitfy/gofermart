package accrual

import (
	"context"

	"github.com/cenkalti/backoff/v5"

	"github.com/spitfy/gofermart/internal/config"
)

type Service struct {
	cfg *config.Config
	s   Storer
}

func NewService(cfg *config.Config, store Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

func (s *Service) Add(ctx context.Context, a Accrual) error {
	add := func() (struct{}, error) {
		err := s.s.Add(ctx, a)
		return struct{}{}, err
	}
	_, err := backoff.Retry(
		ctx,
		add,
		backoff.WithBackOff(backoff.NewExponentialBackOff()), backoff.WithMaxTries(s.cfg.DB.MaxRetries),
	)

	return err
}

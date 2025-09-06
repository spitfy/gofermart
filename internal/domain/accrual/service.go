package accrual

import (
	"context"
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
	return s.s.Add(ctx, a)
}

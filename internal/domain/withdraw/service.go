package withdraw

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

func (s *Service) Add(ctx context.Context, userID int, req Request) error {
	m := Withdraw{
		UserID: userID,
		Order:  req.Order,
		Amount: req.Sum,
	}
	return s.s.Add(ctx, m)
}

func (s *Service) List(ctx context.Context, userID int) ([]Withdraw, error) {
	return s.s.List(ctx, userID)
}

package withdraw

import (
	"context"
	"github.com/spitfy/gofermart/internal/config"
	model "github.com/spitfy/gofermart/internal/model/withdraw"
	repo "github.com/spitfy/gofermart/internal/repository/withdraw"
)

type Service struct {
	cfg *config.Config
	s   repo.Storer
}

func NewService(cfg *config.Config, store repo.Storer) *Service {
	return &Service{
		cfg: cfg,
		s:   store,
	}
}

func (s *Service) Add(ctx context.Context, userID int, req model.Request) error {
	m := model.Withdraw{
		UserID: userID,
		Order:  req.Order,
		Amount: req.Sum,
	}
	return s.s.Add(ctx, m)
}

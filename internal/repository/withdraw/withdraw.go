package withdraw

import (
	"context"
	"github.com/spitfy/gofermart/internal/model/withdraw"
	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	Add(ctx context.Context, w withdraw.Withdraw) error
}

type Store struct {
	*repository.DBStore
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		DBStore: db,
	}
}

func (s *Store) Add(ctx context.Context, w withdraw.Withdraw) error {
	_, err := s.Conn.Exec(ctx,
		`INSERT INTO withdrawals (user_id, order_id, amount) 
					VALUES ($1, (SELECT o.id FROM orders o WHERE o.number = $2), $3)`,
		w.UserID, w.Order, w.Amount,
	)
	return err
}

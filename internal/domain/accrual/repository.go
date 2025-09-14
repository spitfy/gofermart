package accrual

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	Add(ctx context.Context, a Accrual) error
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		pool: db.Pool(),
	}
}

func (s *Store) Add(ctx context.Context, a Accrual) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO accruals (user_id, amount, status, order_id) 
					VALUES ($1, $2, $3, (SELECT o.id FROM orders o WHERE o.number = $4))`,
		a.UserID, a.Amount, a.Status, a.Number,
	)
	return err
}

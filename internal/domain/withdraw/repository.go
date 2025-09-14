package withdraw

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/spitfy/gofermart/internal/repository"
)

var ErrLowBalance = errors.New("not enough balance for withdrawal")

type Storer interface {
	Add(ctx context.Context, w Withdraw) error
	List(ctx context.Context, userID int) ([]Withdraw, error)
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		pool: db.Pool(),
	}
}

func (s *Store) Add(ctx context.Context, w Withdraw) (err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var balance float64
	q := `SELECT balance FROM users WHERE id = $1 FOR UPDATE`
	if err = tx.QueryRow(ctx, q, w.UserID).Scan(&balance); err != nil {
		return err
	}

	if w.Amount > balance {
		return ErrLowBalance
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO withdrawals (user_id, "order", amount) 
					VALUES ($1, $2, $3)`,
		w.UserID, w.Order, w.Amount,
	)

	err = tx.Commit(ctx)
	return err
}

func (s *Store) List(ctx context.Context, userID int) ([]Withdraw, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT w.order, w.amount, w.created_at 
			   FROM withdrawals w
			  WHERE w.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ws []Withdraw
	for rows.Next() {
		var w Withdraw
		if err = rows.Scan(&w.Order, &w.Amount, &w.CreatedAt); err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ws, nil
}

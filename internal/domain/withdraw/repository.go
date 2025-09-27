package withdraw

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/spitfy/gofermart/internal/repository"
)

var ErrLowBalance = errors.New("not enough balance for withdrawal")

//go:generate mockgen -destination=storer_mock.go -package=withdraw github.com/spitfy/gofermart/internal/domain/withdraw Storer
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
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	var balance float64
	q := `SELECT balance FROM users WHERE id = $1 FOR UPDATE`
	if err = tx.QueryRow(ctx, q, w.UserID).Scan(&balance); err != nil {
		return fmt.Errorf("failed to get user balance (userID: %d): %w", w.UserID, err)
	}

	if w.Amount > balance {
		return fmt.Errorf("%w: insufficient balance (available: %.2f, requested: %.2f)",
			ErrLowBalance, balance, w.Amount)
	}
	if _, err = tx.Exec(ctx,
		`INSERT INTO withdrawals (user_id, "order", amount) 
         VALUES ($1, $2, $3)`,
		w.UserID, w.Order, w.Amount,
	); err != nil {
		return fmt.Errorf("failed to insert withdrawal (order: %s): %w", w.Order, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) List(ctx context.Context, userID int) ([]Withdraw, error) {
	const query = `
        SELECT 
            w.order, 
            w.amount, 
            w.created_at 
		FROM withdrawals w
        WHERE w.user_id = $1
        ORDER BY w.created_at DESC`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query withdrawals for user %d: %w", userID, err)
	}
	defer rows.Close()

	var withdrawals []Withdraw
	for rows.Next() {
		var w Withdraw
		if err := rows.Scan(&w.Order, &w.Amount, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan withdrawal row: %w", err)
		}
		withdrawals = append(withdrawals, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after processing withdrawal rows: %w", err)
	}

	return withdrawals, nil
}

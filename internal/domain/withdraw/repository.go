package withdraw

import (
	"context"
	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	Add(ctx context.Context, w Withdraw) error
	List(ctx context.Context, userID int) ([]Withdraw, error)
}

type Store struct {
	*repository.DBStore
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		DBStore: db,
	}
}

func (s *Store) Add(ctx context.Context, w Withdraw) error {
	_, err := s.Conn.Exec(ctx,
		`INSERT INTO withdrawals (user_id, "order", amount) 
					VALUES ($1, $2, $3)`,
		w.UserID, w.Order, w.Amount,
	)
	return err
}

func (s *Store) List(ctx context.Context, userID int) ([]Withdraw, error) {
	rows, err := s.Conn.Query(
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

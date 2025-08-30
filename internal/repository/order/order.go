package order

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/model"
	"github.com/spitfy/gofermart/internal/repository"
)

type Store struct {
	*repository.DBStore
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		DBStore: db,
	}
}

func (s *Store) AddOrder(ctx context.Context, order model.Order) (model.Order, error) {
	_, err := s.Conn.Exec(ctx,
		`INSERT INTO orders (user_id, order_number, status, accrual) VALUES ($1, $2, $3, $4)`,
		order.UserID, order.OrderNumber, order.Status, order.Accrual,
	)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		var userID int
		var status model.OrderStatus
		err = s.Conn.QueryRow(
			ctx,
			"SELECT user_id, status FROM orders WHERE order_number=$1",
			order.OrderNumber,
		).Scan(&userID, &status)
		if err != nil {
			return order, err
		}
		return model.Order{
			UserID: userID,
			Status: status,
		}, nil
	case err != nil:
		return order, err
	default:
		return order, nil
	}
}

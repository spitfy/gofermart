package order

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/model"
	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	AddOrder(ctx context.Context, order model.Order) (model.Order, error)
	ListOrders(ctx context.Context, UserID int) ([]model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

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
		order.UserID, order.Number, order.Status, order.Accrual,
	)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		var userID int
		var status model.OrderStatus
		err = s.Conn.QueryRow(
			ctx,
			"SELECT user_id, status FROM orders WHERE order_number=$1",
			order.Number,
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

func (s *Store) ListOrders(ctx context.Context, userID int) ([]model.Order, error) {
	rows, err := s.Conn.Query(
		ctx,
		"SELECT status, order_number, accrual, created_at FROM orders WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(&o.Status, &o.Number, &o.Accrual, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *Store) Update(ctx context.Context, order model.Order) error {
	_, err := s.Conn.Exec(
		ctx,
		"UPDATE orders SET status = $1, accrual = $2 WHERE order_number = $3",
		order.Status, order.Accrual, order.Number,
	)
	return err
}

package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/repository"
)

//go:generate mockgen -destination=storer_mock.go -package=order github.com/spitfy/gofermart/internal/domain/order Storer
type Storer interface {
	addOrder(ctx context.Context, order Order) (Order, error)
	listOrders(ctx context.Context, UserID int) ([]Order, error)
	listForAccrual(ctx context.Context) ([]orderSend, error)
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		pool: db.Pool(),
	}
}

var ErrUniqueNum = errors.New("order number already exists")

func (s *Store) addOrder(ctx context.Context, order Order) (Order, error) {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO orders (user_id, number) VALUES ($1, $2)`,
		order.UserID, order.Number,
	)

	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		var existing Order
		err = s.pool.QueryRow(
			ctx,
			`SELECT o.user_id, COALESCE(a.status, $1) 
			 FROM orders o  
             LEFT JOIN accruals a ON a.order_id = o.id 
             WHERE number = $2`,
			StatusNew, order.Number,
		).Scan(&existing.UserID, &existing.Status)

		if err != nil {
			return Order{}, fmt.Errorf("fetch order error: %w", err)
		}
		return existing, ErrUniqueNum

	case err != nil:
		return Order{}, fmt.Errorf("failed to create order %s: %w", order.Number, err)
	}
	return order, nil
}

func (s *Store) listOrders(ctx context.Context, userID int) ([]Order, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT coalesce(a.status, $1), o.number, coalesce(a.amount, 0), o.created_at 
			   FROM orders o
					left join accruals a on o.id = a.order_id
			  WHERE o.user_id = $2`,
		StatusNew, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
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

func (s *Store) listForAccrual(ctx context.Context) ([]orderSend, error) {
	rows, err := s.pool.Query(
		ctx,
		`SELECT o.number, o.user_id
			  FROM orders o  
			  LEFT JOIN accruals a ON a.order_id = o.id 
			  WHERE a.order_id IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	var orders []orderSend
	for rows.Next() {
		var o orderSend
		if err = rows.Scan(&o.number, &o.userID); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

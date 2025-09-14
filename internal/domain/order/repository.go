package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	AddOrder(ctx context.Context, order Order) (Order, error)
	ListOrders(ctx context.Context, UserID int) ([]Order, error)
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

func (s *Store) AddOrder(ctx context.Context, order Order) (Order, error) {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO orders (user_id, number) VALUES ($1, $2)`,
		order.UserID, order.Number,
	)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		var userID int
		var st status
		err = s.pool.QueryRow(
			ctx,
			`SELECT o.user_id, coalesce(a.status, $1) 
				   FROM orders o  
				   		left join accruals a on a.order_id = o.id 
				  WHERE number=$2`,
			StatusNew, order.Number,
		).Scan(&userID, &st)
		if err != nil {
			return order, err
		}

		return Order{
			UserID: userID,
			Status: st,
		}, ErrUniqueNum
	case err != nil:
		return order, err
	default:
		return order, nil
	}
}

func (s *Store) ListOrders(ctx context.Context, userID int) ([]Order, error) {
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

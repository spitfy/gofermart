package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	AddOrder(ctx context.Context, order Order) (Order, error)
	ListOrders(ctx context.Context, UserID int) ([]Order, error)
}

type Store struct {
	*repository.DBStore
}

func NewStore(db *repository.DBStore) *Store {
	return &Store{
		DBStore: db,
	}
}

var ErrUniqueNum = errors.New("order number already exists")

func (s *Store) AddOrder(ctx context.Context, order Order) (Order, error) {
	_, err := s.Conn.Exec(ctx,
		`INSERT INTO orders (user_id, number) VALUES ($1, $2)`,
		order.UserID, order.Number,
	)
	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		var userID int
		var status sql.NullString
		err = s.Conn.QueryRow(
			ctx,
			`SELECT o.user_id, a.status 
				   FROM orders o  
				   		left join accruals a on a.order_id = o.id 
				  WHERE number=$1`,
			order.Number,
		).Scan(&userID, &status)
		if err != nil {
			return order, err
		}
		var statusStr string
		if status.Valid {
			statusStr = status.String
		} else {
			statusStr = "" // или любое значение по умолчанию для NULL
		}
		return Order{
			UserID: userID,
			Status: OrderStatus(statusStr),
		}, ErrUniqueNum
	case err != nil:
		return order, err
	default:
		return order, nil
	}
}

func (s *Store) ListOrders(ctx context.Context, userID int) ([]Order, error) {
	rows, err := s.Conn.Query(
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

package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spitfy/gofermart/internal/model"
	"github.com/spitfy/gofermart/internal/repository"
)

type Storer interface {
	AddOrder(ctx context.Context, order model.Order) (model.Order, error)
	ListOrders(ctx context.Context, UserID int) ([]model.Order, error)
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

func (s *Store) AddOrder(ctx context.Context, order model.Order) (model.Order, error) {
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
		return model.Order{
			UserID: userID,
			Status: model.OrderStatus(statusStr),
		}, ErrUniqueNum
	case err != nil:
		return order, err
	default:
		return order, nil
	}
}

func (s *Store) ListOrders(ctx context.Context, userID int) ([]model.Order, error) {
	rows, err := s.Conn.Query(
		ctx,
		`SELECT a.status, o.number, a.amount, o.created_at 
			   FROM orders o
					left join accruals a on o.id = a.order_id
			  WHERE o.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	var status sql.NullString
	var amount sql.NullFloat64
	for rows.Next() {
		var o model.Order
		if err = rows.Scan(&status, &o.Number, &amount, &o.CreatedAt); err != nil {
			return nil, err
		}
		if status.Valid {
			o.Status = model.OrderStatus(status.String)
		} else {
			o.Status = ""
		}
		if amount.Valid {
			o.Accrual = amount.Float64
		} else {
			o.Accrual = 0
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

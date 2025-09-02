package balance

import (
	"context"
	"github.com/spitfy/gofermart/internal/model/balance"
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

func (s *Store) Add(bt balance.BalanceTransaction) error {
	ctx := context.TODO()
	_, err := s.Conn.Exec(ctx,
		`INSERT INTO balance_transactions (user_id, amount, type, order_id) 
					VALUES ($1, $2, $3, (SELECT o.id FROM orders o WHERE o.order_number = $4))`,
		bt.UserID, bt.Amount, bt.Type, bt.OrderNum,
	)
	return err
}

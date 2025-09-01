package balance_transaction

import "time"

type Type string

const (
	TypeAccrual  Type = "ACCRUAL"
	TypeWithdraw Type = "WITHDRAW"
)

type BalanceTransaction struct {
	ID        int       `json:"-"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Type      Type      `json:"type"`
	OrderID   int       `json:"order_id"`
	OrderNum  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

package accrual

import (
	"time"
)

type Accrual struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	OrderID   int       `json:"-"`
	Amount    float64   `json:"accrual"`
	Status    string    `json:"status"`
	Number    string    `json:"number"`
	CreatedAt time.Time `json:"uploaded_at"`
}

type Request struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

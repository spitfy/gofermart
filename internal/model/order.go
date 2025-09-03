package model

import (
	"encoding/json"
	"errors"
	"time"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID          int         `json:"-"`
	UserID      int         `json:"-"`
	OrderNumber string      `json:"number"`
	Status      OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual"`
	CreatedAt   time.Time   `json:"uploaded_at"`
}

func (os OrderStatus) IsValid() bool {
	switch os {
	case StatusNew, StatusProcessing, StatusInvalid, StatusProcessed:
		return true
	default:
		return false
	}
}

func (os *OrderStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	status := OrderStatus(s)
	if !status.IsValid() {
		return errors.New("invalid order status")
	}

	*os = status
	return nil
}

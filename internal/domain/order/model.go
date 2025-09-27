package order

import (
	"encoding/json"
	"errors"
	"time"
)

type status string

const (
	StatusNew        status = "NEW"
	StatusProcessing status = "PROCESSING"
	StatusInvalid    status = "INVALID"
	StatusProcessed  status = "PROCESSED"
)

type Order struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	Number    string    `json:"number"`
	Status    status    `json:"status"`
	Accrual   float64   `json:"accrual"`
	CreatedAt time.Time `json:"uploaded_at"`
}

type orderSend struct {
	userID int
	number string
}

func (os status) IsValid() bool {
	switch os {
	case StatusNew, StatusProcessing, StatusInvalid, StatusProcessed:
		return true
	default:
		return false
	}
}

func (os *status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	status := status(s)
	if !status.IsValid() {
		return errors.New("invalid order status")
	}

	*os = status
	return nil
}

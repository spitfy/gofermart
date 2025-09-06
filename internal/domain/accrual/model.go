package accrual

import (
	"encoding/json"
	"errors"
	"time"
)

type Status string

const (
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
	StatusRegistered Status = "REGISTERED"
)

type Accrual struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	OrderID   int       `json:"-"`
	Amount    float64   `json:"accrual"`
	Status    Status    `json:"status"`
	Number    string    `json:"number"`
	CreatedAt time.Time `json:"uploaded_at"`
}

type Request struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (as Status) IsValid() bool {
	switch as {
	case StatusProcessing, StatusInvalid, StatusProcessed, StatusRegistered:
		return true
	default:
		return false
	}
}

func (as *Status) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	status := Status(s)
	if !status.IsValid() {
		return errors.New("invalid order status")
	}

	*as = status
	return nil
}

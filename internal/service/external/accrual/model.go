package accrual

import (
	"encoding/json"
	"errors"
)

type Response struct {
	Order   string  `json:"order"`
	Status  Status  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type Status string

const (
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
	StatusRegistered Status = "REGISTERED"
)

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

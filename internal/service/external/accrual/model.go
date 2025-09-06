package accrual

import (
	"github.com/spitfy/gofermart/internal/domain/accrual"
)

type Response struct {
	Order   string         `json:"order"`
	Status  accrual.Status `json:"status"`
	Accrual float64        `json:"accrual"`
}

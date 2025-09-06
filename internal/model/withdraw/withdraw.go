package withdraw

import "time"

type Withdraw struct {
	ID        int       `json:"-"`
	UserID    int       `json:"-"`
	Order     string    `json:"order"`
	Amount    float64   `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}

type Request struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

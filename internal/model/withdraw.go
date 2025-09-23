package model

import "time"

type Withdraw struct {
	ID     string `json:"order_id"`
	Number string `json:"order"` // index, unique

	Sum float64 `json:"sum,omitempty"` // nullable

	ProcessedAt time.Time `json:"processed_at"`

	UserID int `json:"user_id"`
}

type AddWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   int64 `json:"sum"`
}

type GetWithdrawalsResponse []*Withdraw

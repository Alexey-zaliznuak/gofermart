package model

import "time"

type Withdraw struct {
	ID     string `json:"order_id"`
	Number string `json:"number"` // index, unique

	Sum float64 `json:"sum,omitempty"` // nullable

	ProcessedAt time.Time `json:"processed_at"`

	UserID int `json:"user_id"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type GetWithdrawalsResponse []WithdrawalResponse

type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"` // RFC3339
}

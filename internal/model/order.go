package model

import "time"

type OrderStatus = string

var (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
)

type Order struct {
	ID string `json:"order_id"`

	Number string      `json:"number"` // index, unique
	Status OrderStatus `json:"status"`

	Accrual *int64 `json:"accrual,omitempty"` // nullable

	UploadedAt time.Time `json:"uploaded_at"`

	UserID int `json:"user_id"`
}

type GetOrdersResponse []*Order

type AccrualResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int   `json:"accrual,omitempty"`
}

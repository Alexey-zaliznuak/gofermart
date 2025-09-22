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

	Accrual *int `json:"accrual,omitempty"` // nullable

	UploadedAt time.Time `json:"uploaded_at"`

	UserID int `json:"user_id"`
}

type GetOrdersResponse []OrderResponse

type OrderResponse struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    *int   `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"` // RFC3339
}

type AccrualResponse struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual *int   `json:"accrual,omitempty"`
}

package models

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	// Идентификатор заказа uuid
	Id          string      `json:"order_uuid"`
	User        string      `json:"user_uuid"`
	Parts       []string    `json:"part_uuids"`
	Total       float64     `json:"total_price"`
	Transaction *string     `json:"transaction_uuid,omitempty"`
	Payment     *string     `json:"payment_method,omitempty"`
	Status      OrderStatus `json:"status"`
}

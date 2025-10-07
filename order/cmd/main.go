package main

import (
	order_v1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	"sync"
	"time"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

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

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*Order),
	}
}

func (s *OrderStorage) getOrder(orderUuid string) *Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) updateOrder(order *Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.Id] = order
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func main() {
	storage := NewOrderStorage()

	OrderHandler := NewOrderHandler(storage)

	orderServer, err := order_v1.NewServer()
}

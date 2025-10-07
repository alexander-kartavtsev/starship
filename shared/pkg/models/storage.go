package models

import "sync"

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

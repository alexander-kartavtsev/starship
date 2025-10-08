package main

import (
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"net/http"
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

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.GetOrderResponse
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.GetOrderResponse),
	}
}

func (s *OrderStorage) getOrder(orderUuid string) *orderV1.GetOrderResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) updateOrder(order *orderV1.GetOrderResponse) {
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

func (h *OrderHandler) CanselOrderById(ctx context.Context, params orderV1.CanselOrderByIdParams) (orderV1.CanselOrderByIdRes, error) {

}

const (
	unknown       orderV1.PaymentMethod = "UNKNOWN"
	card          orderV1.PaymentMethod = "CARD"
	sbp           orderV1.PaymentMethod = "SBP"
	creditCard    orderV1.PaymentMethod = "CREDIT_CARD"
	investorMoney orderV1.PaymentMethod = "INVESTOR_MONEY"
)

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order := &orderV1.GetOrderResponse{
		OrderUUID:       uuid.NewString(),
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      1000,
		TransactionUUID: "",
		PaymentMethod:   "CARD",
	}
}

func (h *OrderHandler) GetOrderByUuid(_ context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order := h.storage.getOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Заказ с id = '" + params.OrderUUID.String() + "' не найден",
		}, nil
	}
	return order, nil
}

func (h *OrderHandler) PayOrderByUuid(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUuidParams) (orderV1.PayOrderByUuidRes, error) {
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusBadGateway,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusBadGateway),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	storage := NewOrderStorage()

	OrderHandler := NewOrderHandler(storage)

	orderServer, err := orderV1.NewServer(OrderHandler)
}

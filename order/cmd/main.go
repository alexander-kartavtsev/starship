package main

import (
	orderV1 "github.com/alexander-kartavtsev/starship/internal/middleware"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"log"
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
	OrderStatusPendingPayment orderV1.OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           orderV1.OrderStatus = "PAID"
	OrderStatusCancelled      orderV1.OrderStatus = "CANCELLED"
)
const (
	unknown       orderV1.PaymentMethod = "UNKNOWN"
	card          orderV1.PaymentMethod = "CARD"
	sbp           orderV1.PaymentMethod = "SBP"
	creditCard    orderV1.PaymentMethod = "CREDIT_CARD"
	investorMoney orderV1.PaymentMethod = "INVESTOR_MONEY"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) getOrder(orderUuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[orderUuid]
	if !ok {
		return nil
	}

	return order
}

func (s *OrderStorage) updateOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (h *OrderHandler) CanselOrderById(_ context.Context, params orderV1.CanselOrderByIdParams) (orderV1.CanselOrderByIdRes, error) {
	order := h.storage.getOrder(params.OrderUUID.String())
	if nil == order {
		return &orderV1.CanselOrderByIdNoContent{
			Code:    http.StatusNotFound,
			Message: "Заказ с uuid = '" + params.OrderUUID.String() + "' не найден",
		}, nil
	}
	if order.Status == OrderStatusPaid {
		return &orderV1.CanselOrderByIdNoContent{
			Code:    http.StatusConflict,
			Message: "Заказ с uuid = '" + params.OrderUUID.String() + "' уже оплачен. Отменить нельзя.",
		}, nil
	}
	if order.Status == OrderStatusPendingPayment {
		order.Status = OrderStatusCancelled
		h.storage.updateOrder(order)
	}
	return &orderV1.CanselOrderByIdNoContent{
		Code:    http.StatusNoContent,
		Message: "Заказ отменен",
	}, nil
}

func (h *OrderHandler) CreateOrder(_ context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	order := &orderV1.OrderDto{
		OrderUUID:       uuid.NewString(),
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      1000,
		TransactionUUID: "",
		PaymentMethod:   unknown,
		Status:          OrderStatusPendingPayment,
	}

	h.storage.updateOrder(order)

	res := &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}

	return res, nil
}

func (h *OrderHandler) GetOrderByUuid(_ context.Context, params orderV1.GetOrderByUuidParams) (orderV1.GetOrderByUuidRes, error) {
	order := h.storage.getOrder(params.OrderUUID.String())
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Заказ с id = '" + params.OrderUUID.String() + "' не найден",
		}, nil
	}
	return order, nil
}

func (h *OrderHandler) PayOrderByUuid(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUuidParams) (orderV1.PayOrderByUuidRes, error) {
	order := h.storage.getOrder(params.OrderUUID.String())

	if order == nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Заказ с id = '" + params.OrderUUID.String() + "' не найден",
		}, nil
	}

	transactionUuid := uuid.NewString()

	order.PaymentMethod = card
	order.Status = OrderStatusPaid
	order.TransactionUUID = transactionUuid

	h.storage.updateOrder(order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUuid,
	}, nil
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
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}
	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

}

package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/alexander-kartavtsev/starship/shared/pkg/middleware"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout      = 5 * time.Second
	shutdownTimeout        = 10 * time.Second
	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
)

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

func (h *OrderHandler) CancelOrderById(_ context.Context, params orderV1.CancelOrderByIdParams) (orderV1.CancelOrderByIdRes, error) {
	order := h.storage.getOrder(params.OrderUUID.String())

	if nil == order {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Заказ с uuid = '" + params.OrderUUID.String() + "' не найден",
		}, nil
	}
	if order.Status == OrderStatusPaid {
		return &orderV1.ConflictError{
			Code:    http.StatusConflict,
			Message: "Заказ с uuid = '" + params.OrderUUID.String() + "' уже оплачен. Отменить нельзя.",
		}, nil
	}
	if order.Status == OrderStatusPendingPayment {
		order.Status = OrderStatusCancelled
		h.storage.updateOrder(order)
	}
	return &orderV1.CancelOrderByIdNoContent{
		Code:    http.StatusNoContent,
		Message: "Заказ отменен",
	}, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	parts, err := getParts(ctx, req.GetPartUuids())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusFailedDependency,
			Message: "При получении данных о запчастях произошла ошибка",
		}, err
	}

	var partUuids []string
	var total float64

	for _, partUuid := range req.GetPartUuids() {
		part, ok := parts[partUuid]
		if !ok {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Деталь с uuid = '" + partUuid + "' не найдена! Уточните заказ",
			}, err
		}
		partUuids = append(partUuids, part.Uuid)
		total += part.Price
	}

	order := &orderV1.OrderDto{
		OrderUUID:  uuid.NewString(),
		UserUUID:   req.UserUUID,
		PartUuids:  partUuids,
		TotalPrice: total,
		Status:     OrderStatusPendingPayment,
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

func (h *OrderHandler) PayOrderByUuid(_ context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderByUuidParams) (orderV1.PayOrderByUuidRes, error) {
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
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}

func getParts(ctx context.Context, partUuids []string) (map[string]*inventoryV1.Part, error) {
	conn, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return nil, err
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	client := inventoryV1.NewInventoryServiceClient(conn)

	inventoryServResp, err := client.ListParts(
		ctx,
		&inventoryV1.ListPartsRequest{
			Filter: &inventoryV1.PartsFilter{
				Uuids: partUuids,
			},
		},
	)
	if err != nil {
		log.Printf("%s\n", err)
		return nil, err
	}
	log.Printf("%v\n", inventoryServResp)

	return inventoryServResp.GetParts(), nil
}

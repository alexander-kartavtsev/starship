package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/net/context"

	v1 "github.com/alexander-kartavtsev/starship/order/internal/api/order/v1"
	orderRepo "github.com/alexander-kartavtsev/starship/order/internal/repository/order"
	orderService "github.com/alexander-kartavtsev/starship/order/internal/service/order"
	customMiddleware "github.com/alexander-kartavtsev/starship/shared/pkg/middleware"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout      = 5 * time.Second
	shutdownTimeout        = 10 * time.Second
	inventoryServerAddress = "localhost:50051"
	paymentServerAddress   = "localhost:50052"
)

func main() {
	repository := orderRepo.NewRepository()
	service := orderService.NewService(repository)
	api := v1.NewApi(service)
	orderServer, err := orderV1.NewServer(api)
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

// getParts - получает из inventory список запчастей по списку partUuids
//func getParts(ctx context.Context, partUuids []string) (map[string]*inventoryV1.Part, error) {
//	conn, err := grpc.NewClient(
//		inventoryServerAddress,
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	if err != nil {
//		log.Printf("failed to connect: %v\n", err)
//		return nil, err
//	}
//	defer func() {
//		if cerr := conn.Close(); cerr != nil {
//			log.Printf("failed to close connect: %v", cerr)
//		}
//	}()
//
//	client := inventoryV1.NewInventoryServiceClient(conn)
//
//	inventoryServResp, err := client.ListParts(
//		ctx,
//		&inventoryV1.ListPartsRequest{
//			Filter: &inventoryV1.PartsFilter{
//				Uuids: partUuids,
//			},
//		},
//	)
//	if err != nil {
//		log.Printf("%s\n", err)
//		return nil, err
//	}
//	log.Printf("%v\n", inventoryServResp)
//
//	return inventoryServResp.GetParts(), nil
//}
//
//func getTransactionUuid(ctx context.Context, orderUuid, userUuid string, paymentMethod paymentV1.PaymentMethod) (string, error) {
//	conn, err := grpc.NewClient(
//		paymentServerAddress,
//		grpc.WithTransportCredentials(insecure.NewCredentials()),
//	)
//	if err != nil {
//		log.Printf("failed to connect: %v\n", err)
//		return "", err
//	}
//	defer func() {
//		if cerr := conn.Close(); cerr != nil {
//			log.Printf("failed to close connect: %v", cerr)
//		}
//	}()
//
//	client := paymentV1.NewPaymentServiceClient(conn)
//
//	paymentServResp, err := client.PayOrder(
//		ctx,
//		&paymentV1.PayOrderRequest{
//			OrderUuid:     orderUuid,
//			UserUuid:      userUuid,
//			PaymentMethod: paymentMethod,
//		},
//	)
//	if err != nil {
//		log.Printf("%s\n", err)
//		return "", err
//	}
//	log.Printf("%v\n", paymentServResp)
//
//	return paymentServResp.GetTransactionUuid(), nil
//}
//
//func convertPaymentMethod(method orderV1.PaymentMethod) paymentV1.PaymentMethod {
//	var paymentMethod paymentV1.PaymentMethod
//	log.Printf("Способ оплаты: %v\n", method)
//	switch method {
//	case card:
//		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
//	case sbp:
//		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
//	case creditCard:
//		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
//	case investorMoney:
//		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
//	default:
//		paymentMethod = paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
//	}
//	log.Printf("Вернулся способ оплаты: %v\n", paymentMethod)
//	return paymentMethod
//}

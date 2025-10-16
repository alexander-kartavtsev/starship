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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/alexander-kartavtsev/starship/order/internal/api/order/v1"
	gRPCinventoryV1 "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/inventory/v1"
	gRPCpaymentV1 "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/alexander-kartavtsev/starship/order/internal/repository/order"
	orderService "github.com/alexander-kartavtsev/starship/order/internal/service/order"
	customMiddleware "github.com/alexander-kartavtsev/starship/shared/pkg/middleware"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
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
	repo := orderRepo.NewRepository()

	connInv := gRPCconn(inventoryServerAddress)
	invClient := gRPCinventoryV1.NewClient(inventoryV1.NewInventoryServiceClient(connInv))
	connPay := gRPCconn(paymentServerAddress)
	payClient := gRPCpaymentV1.NewClient(paymentV1.NewPaymentServiceClient(connPay))

	service := orderService.NewService(repo, invClient, payClient)
	api := v1.NewApi(service)
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	defer func() {
		if cerr := connInv.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
		if cerr := connPay.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

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

func gRPCconn(address string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
	}
	return conn
}

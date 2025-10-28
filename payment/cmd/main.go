package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentApiV1 "github.com/alexander-kartavtsev/starship/payment/internal/api/payment/v1"
	"github.com/alexander-kartavtsev/starship/payment/internal/config"
	paymentService "github.com/alexander-kartavtsev/starship/payment/internal/service/payment"
	"github.com/alexander-kartavtsev/starship/platform/pkg/grpc/health"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

const envPath = "../deploy/compose/payment/.env"

func main() {
	err := config.Load(envPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	conf := config.AppConfig()

	err = logger.Init(conf.Logger.Level(), conf.Logger.AsJson())
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPC.Port()))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	health.RegisterService(s)

	service := paymentService.NewService()
	api := paymentApiV1.NewApi(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("\"🚀 gRPC server listening on %s\n", conf.GRPC.Port())
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}

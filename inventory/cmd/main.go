package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/alexander-kartavtsev/starship/inventory/internal/api/inventory/v1"
	partRepo "github.com/alexander-kartavtsev/starship/inventory/internal/repository/part"
	partService "github.com/alexander-kartavtsev/starship/inventory/internal/service/part"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
	envPath  = "../deploy/compose/inventory/.env"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	// Проверяем соединение с базой данных
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database: %v\n", err)
		return
	}

	// Получаем базу данных
	db := client.Database("inventory")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
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

	repo := partRepo.NewRepository(db)
	service := partService.NewService(repo)
	api := v1.NewApi(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("\"🚀 gRPC server listening on %d\n", grpcPort)
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

package main

import (
	"context"
	"fmt"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// *********** server
const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	parts map[string]*inventoryV1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	s.parts["bc838184-47b8-8086-1e05-3060f2c0ada1"] = &inventoryV1.Part{
		Uuid:          "bc838184-47b8-8086-1e05-3060f2c0ada1",
		Name:          "Тестовое крыло",
		Description:   "Для полетов во сне и наяву",
		Price:         100000000,
		StockQuantity: 10,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Dimensions: &inventoryV1.Dimensions{
			Weight: 120,
			Length: 24,
			Height: 7,
			Width:  12,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Angel",
			Country: "Россия",
			Website: "russia-angel.ru",
		},
		Tags:      []string{"крыло"},
		Metadata:  map[string]*inventoryV1.Value{},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	part, ok := s.parts[req.GetUuid()]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Запчасть с UUID %s не найдена", req.GetUuid())
	}
	return &inventoryV1.GetPartResponse{
		Info: part,
	}, nil
}

//func (s *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartRequest) (*inventoryV1.ListPartResponse, error) {
//partFilter := req.GetFilter()
//}

//******************* client

const serverAddress = "localhost:50051"

//******************* main

func main() {
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

	service := &inventoryService{
		parts: make(map[string]*inventoryV1.Part),
	}

	inventoryV1.RegisterInventoryServiceServer(s, service)

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

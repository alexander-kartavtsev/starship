package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

const serverAddress = "localhost:50051"

func main() {
	ctx := context.Background()

	conn, err := grpc.NewClient(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	client := inventoryV1.NewInventoryServiceClient(conn)

	resp, err := client.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: "bc838184-47b8-8086-1e05-3060f2c0ada1"})
	if err != nil {
		log.Printf("%s\n", err)
	}
	log.Printf("%v\n", resp)
}

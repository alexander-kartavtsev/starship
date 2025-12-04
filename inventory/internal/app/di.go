package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiV1 "github.com/alexander-kartavtsev/starship/inventory/internal/api/inventory/v1"
	"github.com/alexander-kartavtsev/starship/inventory/internal/config"
	"github.com/alexander-kartavtsev/starship/inventory/internal/repository"
	partRepository "github.com/alexander-kartavtsev/starship/inventory/internal/repository/part"
	"github.com/alexander-kartavtsev/starship/inventory/internal/service"
	partService "github.com/alexander-kartavtsev/starship/inventory/internal/service/part"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	"github.com/alexander-kartavtsev/starship/platform/pkg/tracing"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryApi        inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository
	mongoDb             *mongo.Database
	mongoClient         *mongo.Client
	authServiceClient   authV1.AuthServiceClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryApi(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryApi == nil {
		d.inventoryApi = apiV1.NewApi(d.InventoryService(ctx))
		logger.Info(ctx, "Инициализация Api")
	}
	return d.inventoryApi
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = partService.NewService(d.InventoryRepository(ctx))
		logger.Info(ctx, "Инициализация Service")
	}
	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = partRepository.NewRepository(ctx, d.MongoDb(ctx))
		logger.Info(ctx, "Инициализация Repository")
	}
	return d.inventoryRepository
}

func (d *diContainer) MongoDb(ctx context.Context) *mongo.Database {
	if d.mongoDb == nil {
		d.mongoDb = d.MongoClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
		logger.Info(ctx, "Инициализация MongoDb")
	}
	return d.mongoDb
}

func (d *diContainer) MongoClient(ctx context.Context) *mongo.Client {
	if d.mongoClient == nil {
		client, err := mongo.Connect(
			ctx,
			options.Client().ApplyURI(config.AppConfig().Mongo.Uri()),
			options.Client().SetServerSelectionTimeout(5*time.Second),
		)
		if err != nil {
			logger.Error(ctx, "Ошибка при создании клиента", zap.Error(err))
			panic(fmt.Sprintf("failed to create MongoClient: %s\n", err.Error()))
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			logger.Error(ctx, "Ошибка при проверке ping", zap.Error(err))
			panic(fmt.Sprintf("failed to ping MongoClient: %v\n", err))
		}
		closer.AddNamed("MongoDb client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		logger.Info(ctx, "Инициализация MongoClient")

		d.mongoClient = client
	}
	return d.mongoClient
}

func (d *diContainer) AuthClient(ctx context.Context) authV1.AuthServiceClient {
	if d.authServiceClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IamClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("iam-service")),
		)
		if err != nil {
			log.Printf("failed to connect: %v\n", err)
		}

		logger.Info(ctx, "Инициализация AuthService gRPC Client")

		closer.AddNamed("AuthService gRPC client", func(ctx context.Context) error {
			return conn.Close()
		})

		d.authServiceClient = authV1.NewAuthServiceClient(conn)
	}
	return d.authServiceClient
}

package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/alexander-kartavtsev/starship/order/internal/api/order/v1"
	orderGRPC "github.com/alexander-kartavtsev/starship/order/internal/client/grpc"
	gRPCinventoryV1 "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/inventory/v1"
	gRPCpaymentV1 "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/payment/v1"
	"github.com/alexander-kartavtsev/starship/order/internal/config"
	"github.com/alexander-kartavtsev/starship/order/internal/repository"
	orderRepo "github.com/alexander-kartavtsev/starship/order/internal/repository/order"
	"github.com/alexander-kartavtsev/starship/order/internal/service"
	"github.com/alexander-kartavtsev/starship/order/internal/service/order"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderServer     *orderV1.Server
	orderApi        orderV1.Handler
	orderService    service.OrderService
	invClient       orderGRPC.InventoryClient
	payClient       orderGRPC.PaymentClient
	orderRepository repository.OrderRepository
	dbConn          *pgx.Conn
	dbPool          *pgxpool.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderServer(ctx context.Context) *orderV1.Server {
	if d.orderServer == nil {
		orderServer, err := orderV1.NewServer(d.OrderApi(ctx))
		if err != nil {
			logger.Error(ctx, "Ошибка инициализации сервера", zap.Error(err))
		}
		logger.Info(ctx, "Инициализация сервера")
		d.orderServer = orderServer
	}
	return d.orderServer
}

func (d *diContainer) OrderApi(ctx context.Context) orderV1.Handler {
	if d.orderApi == nil {
		d.orderApi = v1.NewApi(d.OrderService(ctx))
		logger.Info(ctx, "Инициализация Api")
	}
	return d.orderApi
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = order.NewService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
		logger.Info(ctx, "Инициализация сервиса")
	}
	return d.orderService
}

func (d *diContainer) InventoryClient(ctx context.Context) orderGRPC.InventoryClient {
	if d.invClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("failed to connect: %v\n", err)
		}
		logger.Info(ctx, "Инициализация Inventory")

		closer.AddNamed("Inventory gRPC client", func(ctx context.Context) error {
			return conn.Close()
		})

		d.invClient = gRPCinventoryV1.NewClient(inventoryV1.NewInventoryServiceClient(conn))
	}
	return d.invClient
}

func (d *diContainer) PaymentClient(ctx context.Context) orderGRPC.PaymentClient {
	if d.payClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Printf("failed to connect: %v\n", err)
		}

		logger.Info(ctx, "Инициализация Payment")

		closer.AddNamed("Payment gRPC client", func(ctx context.Context) error {
			return conn.Close()
		})

		d.payClient = gRPCpaymentV1.NewClient(paymentV1.NewPaymentServiceClient(conn))
	}
	return d.payClient
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepo.NewRepository(d.DbConn(ctx), d.DbPool(ctx))
		logger.Info(ctx, "Инициализация репозитория")
	}
	return d.orderRepository
}

func (d *diContainer) DbConn(ctx context.Context) *pgx.Conn {
	if d.dbConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.Uri())
		if err != nil {
			log.Printf("failed to connect to database: %v\n", err)
			panic(err)
		}

		// Проверяем, что соединение с базой установлено
		err = conn.Ping(ctx)
		if err != nil {
			log.Printf("База данных недоступна: %v\n", err)
			panic(err)
		}

		logger.Info(ctx, "Инициализация DbConn")

		closer.AddNamed("Db connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		d.dbConn = conn
	}
	return d.dbConn
}

func (d *diContainer) DbPool(ctx context.Context) *pgxpool.Pool {
	if d.dbPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.Uri())
		if err != nil {
			log.Printf("Ошибка подключения к б/д: %v\n", err)
		}

		logger.Info(ctx, "Инициализация DbPool")

		closer.AddNamed("Db pool connection", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.dbPool = pool
	}
	return d.dbPool
}

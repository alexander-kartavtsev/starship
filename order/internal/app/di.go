package app

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/alexander-kartavtsev/starship/order/internal/api/order/v1"
	orderGRPC "github.com/alexander-kartavtsev/starship/order/internal/client/grpc"
	gRPCinventoryV1 "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/inventory/v1"
	gRPCpaymentV1 "github.com/alexander-kartavtsev/starship/order/internal/client/grpc/payment/v1"
	"github.com/alexander-kartavtsev/starship/order/internal/config"
	kafkaConverter "github.com/alexander-kartavtsev/starship/order/internal/converter/kafka"
	"github.com/alexander-kartavtsev/starship/order/internal/converter/kafka/decoder"
	"github.com/alexander-kartavtsev/starship/order/internal/repository"
	orderRepo "github.com/alexander-kartavtsev/starship/order/internal/repository/order"
	"github.com/alexander-kartavtsev/starship/order/internal/service"
	orderConsumer "github.com/alexander-kartavtsev/starship/order/internal/service/consumer/order_consumer"
	"github.com/alexander-kartavtsev/starship/order/internal/service/order"
	orderProducer "github.com/alexander-kartavtsev/starship/order/internal/service/producer/order_producer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	wrappedKafka "github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/alexander-kartavtsev/starship/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/alexander-kartavtsev/starship/platform/pkg/kafka/producer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	kafkaMiddleware "github.com/alexander-kartavtsev/starship/platform/pkg/middleware/kafka"
	migrator "github.com/alexander-kartavtsev/starship/platform/pkg/migrator/pg"
	orderV1 "github.com/alexander-kartavtsev/starship/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderServer  *orderV1.Server
	orderApi     orderV1.Handler
	orderService service.OrderService

	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	invClient       orderGRPC.InventoryClient
	payClient       orderGRPC.PaymentClient
	orderRepository repository.OrderRepository

	consumerGroup sarama.ConsumerGroup
	orderConsumer wrappedKafka.Consumer
	orderDecoder  kafkaConverter.OrderDecoder
	syncProducer  sarama.SyncProducer
	orderProducer wrappedKafka.Producer

	dbConn   *pgx.Conn
	dbPool   *pgxpool.Pool
	migrator *migrator.Migrator
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
		logger.Info(ctx, "Инициализация Http Server")
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
		d.orderService = order.NewService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx), d.OrderProducerService(ctx))
		logger.Info(ctx, "Инициализация Service")
	}
	return d.orderService
}

func (d *diContainer) OrderProducerService(ctx context.Context) service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderProducer(ctx))
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderConsumer(ctx), d.OrderDecoder(), d.OrderRepository(ctx))
	}

	return d.orderConsumerService
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
		logger.Info(ctx, "Инициализация Inventory gRPC Client")

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

		logger.Info(ctx, "Инициализация Payment gRPC Client")

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
		logger.Info(ctx, "Инициализация Repository")
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

func (d *diContainer) Migrator(ctx context.Context) *migrator.Migrator {
	if d.migrator == nil {
		d.migrator = migrator.NewMigrator(
			stdlib.OpenDB(*d.DbConn(ctx).Config().Copy()),
			config.AppConfig().Postgres.MigrationsDir(),
		)
	}
	return d.migrator
}

func (d *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderConsumer.GroupID(),
			config.AppConfig().OrderConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}

		logger.Info(ctx, "Инициализация ConsumerGroup")

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderConsumer == nil {
		d.orderConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(ctx),
			[]string{
				config.AppConfig().OrderConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderConsumer
}

func (d *diContainer) OrderDecoder() kafkaConverter.OrderDecoder {
	if d.orderDecoder == nil {
		d.orderDecoder = decoder.NewAssemblyDecoder()
	}

	return d.orderDecoder
}

func (d *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("Ошибка инициализации sync producer: %s\n", err.Error()))
		}

		logger.Info(ctx, "Инициализация SyncProducer")

		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}
	return d.syncProducer
}

func (d *diContainer) OrderProducer(ctx context.Context) wrappedKafka.Producer {
	if d.orderProducer == nil {
		d.orderProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig().OrderProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.orderProducer
}

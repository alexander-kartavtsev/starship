package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/alexander-kartavtsev/starship/assembly/internal/config"
	kafkaConverter "github.com/alexander-kartavtsev/starship/assembly/internal/converter/kafka"
	"github.com/alexander-kartavtsev/starship/assembly/internal/converter/kafka/decoder"
	"github.com/alexander-kartavtsev/starship/assembly/internal/service"
	assemblyConsumer "github.com/alexander-kartavtsev/starship/assembly/internal/service/consumer/assembly_consumer"
	assemblyProducer "github.com/alexander-kartavtsev/starship/assembly/internal/service/producer/assembly_producer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	wrappedKafka "github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/alexander-kartavtsev/starship/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/alexander-kartavtsev/starship/platform/pkg/kafka/producer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	kafkaMiddleware "github.com/alexander-kartavtsev/starship/platform/pkg/middleware/kafka"
)

type diContainer struct {
	producerService service.ProducerService
	consumerService service.ConsumerService

	consumerGroup sarama.ConsumerGroup
	consumer      wrappedKafka.Consumer
	decoder       kafkaConverter.AssemblyDecoder
	syncProducer  sarama.SyncProducer
	producer      wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ProducerService(ctx context.Context) service.ProducerService {
	if d.producerService == nil {
		d.producerService = assemblyProducer.NewService(d.Producer(ctx))
	}

	return d.producerService
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = assemblyConsumer.NewService(d.Consumer(ctx), d.Decoder(), d.ProducerService(ctx))
	}

	return d.consumerService
}

func (d *diContainer) ConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyConsumer.GroupID(),
			config.AppConfig().AssemblyConsumer.Config(),
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

func (d *diContainer) Consumer(ctx context.Context) wrappedKafka.Consumer {
	if d.consumer == nil {
		d.consumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(ctx),
			[]string{
				config.AppConfig().AssemblyConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.consumer
}

func (d *diContainer) Decoder() kafkaConverter.AssemblyDecoder {
	if d.decoder == nil {
		d.decoder = decoder.NewAssemblyDecoder()
	}

	return d.decoder
}

func (d *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyProducer.Config(),
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

func (d *diContainer) Producer(ctx context.Context) wrappedKafka.Producer {
	if d.producer == nil {
		d.producer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig().AssemblyProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.producer
}

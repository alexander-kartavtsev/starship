package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/alexander-kartavtsev/starship/notification/internal/client/http"
	telegramClient "github.com/alexander-kartavtsev/starship/notification/internal/client/http/telegram"
	"github.com/alexander-kartavtsev/starship/notification/internal/config"
	kafkaConverter "github.com/alexander-kartavtsev/starship/notification/internal/converter/kafka"
	"github.com/alexander-kartavtsev/starship/notification/internal/converter/kafka/decoder"
	"github.com/alexander-kartavtsev/starship/notification/internal/service"
	orderAssembledConsumer "github.com/alexander-kartavtsev/starship/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/alexander-kartavtsev/starship/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/alexander-kartavtsev/starship/notification/internal/service/telegram"
	"github.com/alexander-kartavtsev/starship/platform/pkg/closer"
	wrappedKafka "github.com/alexander-kartavtsev/starship/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/alexander-kartavtsev/starship/platform/pkg/kafka/consumer"
	"github.com/alexander-kartavtsev/starship/platform/pkg/logger"
	kafkaMiddleware "github.com/alexander-kartavtsev/starship/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderAssembledConsumerService service.OrderAssembledConsumerService
	orderPaidConsumerService      service.OrderPaidConsumerService
	telegramService               service.TelegramService

	telegramClient httpClient.TelegramClient
	telegramBot    *bot.Bot

	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderPaidConsumerGroup      sarama.ConsumerGroup

	orderAssembledConsumer wrappedKafka.Consumer
	orderPaidConsumer      wrappedKafka.Consumer

	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder
	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.OrderAssembledConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(d.OrderAssembledConsumer(ctx), d.OrderAssemblyDecoder(), d.TelegramService(ctx))
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(d.OrderPaidConsumer(ctx), d.OrderPaidDecoder(), d.TelegramService(ctx))
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}

		logger.Info(ctx, "Инициализация OrderAssembledConsumerGroup")

		closer.AddNamed("Kafka OrderAssembledConsumerGroup", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderPaidConsumerGroup(ctx context.Context) sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}

		logger.Info(ctx, "Инициализация OrderPaidConsumerGroup")

		closer.AddNamed("Kafka OrderPaidConsumerGroup", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderAssembledConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(ctx),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(ctx),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssemblyDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}

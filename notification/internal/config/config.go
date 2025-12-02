package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alexander-kartavtsev/starship/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledConsumer ConsumerConfig
	OrderPaidConsumer      ConsumerConfig
	TelegramBot            TelegramBotConfig
	Tracing                TracingConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumer, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumer, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	telegramBot, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	tracingConfig, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerConfig,
		Kafka:                  kafkaConfig,
		OrderAssembledConsumer: orderAssembledConsumer,
		OrderPaidConsumer:      orderPaidConsumer,
		TelegramBot:            telegramBot,
		Tracing:                tracingConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}

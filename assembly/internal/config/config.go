package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alexander-kartavtsev/starship/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	Logger           LoggerConfig
	Kafka            KafkaConfig
	AssemblyProducer AssemblyProducerConfig
	AssemblyConsumer AssemblyConsumerConfig
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

	assemblyConsumerConfig, err := env.NewAssemblyConsumerConfig()
	if err != nil {
		return err
	}

	assemblyProducerConfig, err := env.NewAssemblyProducerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:           loggerConfig,
		Kafka:            kafkaConfig,
		AssemblyConsumer: assemblyConsumerConfig,
		AssemblyProducer: assemblyProducerConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}

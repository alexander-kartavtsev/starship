package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderProducerEnvConfig struct {
	TopicName string `env:"ORDER_PAID_TOPIC_NAME,required"`
}

type orderProducerConfig struct {
	raw orderProducerEnvConfig
}

func NewOrderProducerConfig() (*orderProducerConfig, error) {
	var raw orderProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderProducerConfig{raw: raw}, nil
}

func (cfg *orderProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *orderProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}

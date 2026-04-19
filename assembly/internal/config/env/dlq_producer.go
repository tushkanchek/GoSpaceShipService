package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type dlqProducerEnvConfig struct {
	TopicName string `env:"DLQ_TOPIC_NAME,required"`
}

type dlqProducerConfig struct {
	raw dlqProducerEnvConfig
}

func NewDLQProducerConfig() (*dlqProducerConfig, error) {
	var raw dlqProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &dlqProducerConfig{
		raw: raw,
	}, nil
}

func (cfg *dlqProducerConfig) Topic() string {
	return cfg.raw.TopicName
}

func (cfg *dlqProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}

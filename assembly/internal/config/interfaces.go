package config

import (
	"time"

	"github.com/IBM/sarama"
)

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type DLQProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type RetryConfig interface {
	MaxRetries() int
	InitialDelay() time.Duration
	MaxDelay() time.Duration
	BackoffMultiplier() float64
}

type TestingConfig interface {
	FailProduceForTesting() bool
}

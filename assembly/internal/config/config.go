package config

import (
	"assembly/internal/config/env"
	"github.com/joho/godotenv"
)

var AppConfig *config

type config struct {
	Logger            LoggerConfig
	Kafka             KafkaConfig
	OrderPaidConsumer OrderPaidConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	AppConfig = &config{
		Logger:            loggerCfg,
		Kafka:             kafkaCfg,
		OrderPaidConsumer: orderPaidConsumerCfg,
	}

	return nil
}

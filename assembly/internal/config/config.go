package config

import (
	"assembly/internal/config/env"

	"github.com/joho/godotenv"
)

var AppConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	OrderAssembledProducer OrderAssembledProducerConfig
	DLQProducer            DLQProducerConfig
	Retry                  RetryConfig
	Testing                TestingConfig
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

	orderAssembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	dlqProducerCfg, err := env.NewDLQProducerConfig()
	if err != nil {
		return err
	}

	retryCfg, err := env.NewRetryConfig()
	if err != nil {
		return err
	}

	testingCfg := env.LoadTestingConfig()

	AppConfig = &config{
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		OrderAssembledProducer: orderAssembledProducerCfg,
		DLQProducer:            dlqProducerCfg,
		Retry:                  retryCfg,
		Testing:                testingCfg,
	}

	return nil
}

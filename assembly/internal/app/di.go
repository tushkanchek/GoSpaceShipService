package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"assembly/internal/config"
	kafkaConverter "assembly/internal/converter/kafka"
	"assembly/internal/converter/kafka/decoder"
	"assembly/internal/service"
	"assembly/internal/service/consumer/order_consumer"
	wrappedKafka "platform/pkg/kafka"
	wrappedKafkaConsumer "platform/pkg/kafka/consumer"
	"platform/pkg/closer"
	"platform/pkg/logger"
)

type diContainer struct {
	consumerGroup    sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
	consumerService  service.ConsumerService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) ConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig.Kafka.Brokers(),
			config.AppConfig.OrderPaidConsumer.GroupID(),
			config.AppConfig.OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Errorf("failed to create consumer group: %w", err))
		}

		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}
	return d.consumerGroup
}

func (d *diContainer) OrderPaidConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(ctx),
			[]string{config.AppConfig.OrderPaidConsumer.Topic()},
			logger.Logger(),
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}

func (d *diContainer) ConsumerService(ctx context.Context) service.ConsumerService {
	if d.consumerService == nil {
		d.consumerService = order_consumer.NewService(
			d.OrderPaidConsumer(ctx),
			d.OrderPaidDecoder(),
		)
	}
	return d.consumerService
}

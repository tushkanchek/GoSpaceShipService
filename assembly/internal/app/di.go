package app

import (
	"context"
	"fmt"

	"assembly/internal/config"
	kafkaConverter "assembly/internal/converter/kafka"
	"assembly/internal/converter/kafka/decoder"
	"assembly/internal/converter/kafka/encoder"
	"assembly/internal/service"
	"assembly/internal/service/consumer/order_consumer"
	orderAssembledProducer "assembly/internal/service/producer/order_assembled_producer"
	"github.com/IBM/sarama"
	"platform/pkg/closer"
	wrappedKafka "platform/pkg/kafka"
	wrappedKafkaConsumer "platform/pkg/kafka/consumer"
	wrappedKafkaProducer "platform/pkg/kafka/producer"
	"platform/pkg/logger"
)

type diContainer struct {
	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder
	consumerService   service.AssemblyConsumerService
	producerService   service.AssemblyProducerService

	syncProducer           sarama.SyncProducer
	orderAssembledProducer wrappedKafka.Producer
	orderAssembledEncoder  kafkaConverter.OrderAssembledEncoder
	dlqProducer            wrappedKafka.Producer
	dlqEncoder             kafkaConverter.DLQEventEncoder
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

func (d *diContainer) ProducerService(ctx context.Context) service.AssemblyProducerService {
	if d.producerService == nil {
		d.producerService = orderAssembledProducer.NewService(
			d.OrderAssembledProducer(ctx),
			d.OrderAssembledEncoder(),
		)
	}
	return d.producerService
}

func (d *diContainer) ConsumerService(ctx context.Context) service.AssemblyConsumerService {
	if d.consumerService == nil {
		d.consumerService = order_consumer.NewService(
			d.OrderPaidConsumer(ctx),
			d.OrderPaidDecoder(),
			d.ProducerService(ctx),
			d.DLQProducer(ctx),
			d.DLQEncoder(),
		)
	}
	return d.consumerService
}

func (d *diContainer) SyncProducer(ctx context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig.Kafka.Brokers(),
			config.AppConfig.OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka Sync order assembled producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssembledProducer(ctx context.Context) wrappedKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig.OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderAssembledProducer
}

func (d *diContainer) OrderAssembledEncoder() kafkaConverter.OrderAssembledEncoder {
	if d.orderAssembledEncoder == nil {
		d.orderAssembledEncoder = encoder.NewOrderAssembledEncoder()
	}
	return d.orderAssembledEncoder
}

func (d *diContainer) DLQProducer(ctx context.Context) wrappedKafka.Producer {
	if d.dlqProducer == nil {
		d.dlqProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig.DLQProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.dlqProducer
}

func (d *diContainer) DLQEncoder() kafkaConverter.DLQEventEncoder {
	if d.dlqEncoder == nil {
		d.dlqEncoder = encoder.NewDLQEventEncoder()
	}
	return d.dlqEncoder
}

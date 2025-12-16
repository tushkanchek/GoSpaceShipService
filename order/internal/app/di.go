package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderV1API "order/internal/api/order/v1"
	grpcClients "order/internal/client/grpc"
	clientInventory "order/internal/client/grpc/inventory/v1"
	clientPayment "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/repository"
	orderRepository "order/internal/repository/order"
	"order/internal/service"
	orderService "order/internal/service/order"
	orderproducer "order/internal/service/producer/order_producer"
	"platform/pkg/closer"
	wrappedKafka "platform/pkg/kafka"
	wrappedKafkaProducer "platform/pkg/kafka/producer"
	"platform/pkg/logger"
	orderV1 "shared/pkg/openapi/order/v1"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService         service.OrderService
	orderProducerService service.OrderProducerService

	orderRepository repository.OrderRepository

	inventoryClient grpcClients.InventoryClient

	paymentClient grpcClients.PaymentClient

	postgresPool *pgxpool.Pool

	inventoryGRPCConn *grpc.ClientConn

	paymentGRPCConn *grpc.ClientConn

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
			d.OrderProducerService(),
		)
	}
	return d.orderService
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderproducer.NewService(d.OrderPaidProducer())
	}
	return d.orderProducerService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewOrderRepository(d.PostgresPool(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) grpcClients.InventoryClient {
	if d.inventoryClient == nil {
		conn := d.InventoryGRPCConn(ctx)
		generatedClient := inventoryV1.NewInventoryServiceClient(conn)
		d.inventoryClient = clientInventory.NewClient(generatedClient)
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) grpcClients.PaymentClient {
	if d.paymentClient == nil {
		conn := d.PaymentGRPCConn(ctx)
		generatedClient := paymentV1.NewPaymentServiceClient(conn)
		d.paymentClient = clientPayment.NewClient(generatedClient)
	}
	return d.paymentClient
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Errorf("failed to create postgres connection pool: %w", err))
		}

		err = pool.Ping(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to ping postgres: %w", err))
		}

		closer.AddNamed("PostgreSQL pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.postgresPool = pool
	}
	return d.postgresPool
}

func (d *diContainer) InventoryGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.inventoryGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Adress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to inventory gRPC: %w", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryGRPCConn = conn
	}
	return d.inventoryGRPCConn
}

func (d *diContainer) PaymentGRPCConn(ctx context.Context) *grpc.ClientConn {
	if d.paymentGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Adress(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to connect to payment gRPC: %w", err))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentGRPCConn = conn
	}
	return d.paymentGRPCConn
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka Sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderPaidProducer() wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
}

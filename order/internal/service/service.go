package service

import (
	"context"

	"github.com/google/uuid"
	"order/internal/model"
)

type OrderService interface {
	GetOrderByUUID(ctx context.Context, order_uuid uuid.UUID) (*model.Order, error)
	CreateOrder(ctx context.Context, user_uuid uuid.UUID, part_uuids []uuid.UUID) (*model.Order, error)
	PayOrder(ctx context.Context, order_uuid uuid.UUID, PaymentMethod model.PaymentMethod) (uuid.UUID, error)
	CancelOrder(ctx context.Context, order_uuid uuid.UUID) error
	UpdateOrderStatus(ctx context.Context, order_uuid uuid.UUID, status model.OrderStatus) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error
}

type OrderConsumerService interface {
	RunConsumer(ctx context.Context) error
}
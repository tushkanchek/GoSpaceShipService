package repository

import (
	"context"

	"github.com/google/uuid"
	model "order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, order_uuid uuid.UUID) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}

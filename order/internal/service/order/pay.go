package order

import (
	"context"

	"order/internal/model"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// TODO: check orderstatus cancel
func (s *service) PayOrder(ctx context.Context, order_uuid uuid.UUID, PaymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if order_uuid == uuid.Nil {
		return uuid.Nil, model.ErrEmptyOrderUuid
	}
	order, err := s.OrderRepository.GetOrder(ctx, order_uuid)
	if err != nil {
		return uuid.Nil, err
	}
	if order == nil {
		return uuid.Nil, model.ErrOrderNotFound
	}
	if order.OrderStatus == model.OrderStatusPAID {
		return uuid.Nil, model.ErrPayOrderStatusPaid
	}
	if order.OrderStatus == model.OrderStatusCANCELLED {
		return uuid.Nil, model.ErrPayOrderStatusCancelled
	}

	transaction_uuid, err := s.PaymentClient.PayOrder(ctx, order_uuid.String(), order.UserUUID.String(), PaymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	order.OrderStatus = model.OrderStatusPAID
	order.TransactionUUID = lo.ToPtr(transaction_uuid)
	order.PaymentMethod = &PaymentMethod

	err = s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}

	return transaction_uuid, err
}

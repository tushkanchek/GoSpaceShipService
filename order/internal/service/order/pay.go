package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"order/internal/model"
)

// TODO: check orderstatus cancel
func (s *service) PayOrder(ctx context.Context, order_uuid uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
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

	transaction_uuid, err := s.PaymentClient.PayOrder(ctx, order_uuid.String(), order.UserUUID.String(), paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	order.OrderStatus = model.OrderStatusPAID
	order.TransactionUUID = lo.ToPtr(transaction_uuid)
	order.PaymentMethod = &paymentMethod

	err = s.OrderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}

	return transaction_uuid, err
}

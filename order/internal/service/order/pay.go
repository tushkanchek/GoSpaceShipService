package order

import (
	"context"
	"platform/pkg/logger"

	"order/internal/model"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// TODO: check orderstatus cancel
func (s *service) PayOrder(ctx context.Context, order_uuid uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if order_uuid == uuid.Nil {
		return uuid.Nil, model.ErrEmptyOrderUuid
	}

	reqGetCtx, cancelGet := context.WithTimeout(ctx, model.RequestTimeOutRead)
	defer cancelGet()

	order, err := s.OrderRepository.GetOrder(reqGetCtx, order_uuid)
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

	reqPayCtx, cancelPay := context.WithTimeout(ctx, model.RequestTimeOutUpdate)
	defer cancelPay()

	transaction_uuid, err := s.PaymentClient.PayOrder(reqPayCtx, order_uuid.String(), order.UserUUID.String(), paymentMethod)
	if transaction_uuid == uuid.Nil {
		return uuid.Nil, err
	}
	if err != nil {
		return uuid.Nil, err
	}

	order.OrderStatus = model.OrderStatusPAID
	order.TransactionUUID = lo.ToPtr(transaction_uuid)
	order.PaymentMethod = &paymentMethod

	reqUpdateCtx, cancelUpdate := context.WithTimeout(ctx, model.RequestTimeOutUpdate)
	defer cancelUpdate()

	err = s.OrderRepository.UpdateOrder(reqUpdateCtx, order)
	if err != nil {
		return uuid.Nil, err
	}

	eventUuid, err := uuid.NewUUID()
	if err != nil {
		logger.Error(ctx, "couldn't create event uuid", zap.Error(err))
		return uuid.Nil, err
	}
	//Kafka Producer
	err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
		EventUuid:     eventUuid,
		OrderUuid:     order_uuid,
		UserUuid:      order.UserUUID,
		PaymentMethod: int32(paymentMethod),
		TransactionUuid: transaction_uuid,
	})
	if err != nil {
		logger.Error(ctx, "couldn't produce event transaction uuid")
		return transaction_uuid, err
	}

	return transaction_uuid, err
}

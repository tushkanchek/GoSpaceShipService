package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSucces() {
	orderUuid := uuid.MustParse(gofakeit.UUID())

	userUuid := uuid.MustParse(gofakeit.UUID())

	paymentMethod := model.PaymentMethodCARD

	order := &model.Order{
		OrderUUID:   orderUuid,
		UserUUID:    userUuid,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	transaction_uuid := uuid.MustParse(gofakeit.UUID())

	s.orderRepo.On("GetOrder", mock.Anything, orderUuid).
		Return(order, nil).
		Once()

	s.paymentClient.
		On("PayOrder", mock.Anything, orderUuid.String(), userUuid.String(), paymentMethod).
		Return(transaction_uuid, nil).
		Once()

	s.orderRepo.
		On("UpdateOrder", mock.Anything, mock.MatchedBy(func(o *model.Order) bool {
			return o.OrderUUID == orderUuid &&
				o.OrderStatus == model.OrderStatusPAID &&
				o.TransactionUUID != nil &&
				*o.TransactionUUID == transaction_uuid
		})).
		Return(nil).
		Once()

	s.orderPaidProducer.
		On("ProduceOrderPaid", mock.Anything, mock.MatchedBy(func(e model.OrderPaidEvent) bool {
			return e.OrderUuid == orderUuid &&
				e.TransactionUuid == transaction_uuid
		})).
		Return(nil).
		Once()

	result, err := s.service.PayOrder(context.Background(), orderUuid, paymentMethod)

	s.NoError(err)
	s.Equal(transaction_uuid, result)
}

func (s *ServiceSuite) TestPayOrderEmptyOrderUuid() {
	orderUuid := uuid.Nil

	paymentMethod := model.PaymentMethodCARD

	ctx := context.Background()
	result, err := s.service.PayOrder(ctx, orderUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())
	s.Empty(result)
}

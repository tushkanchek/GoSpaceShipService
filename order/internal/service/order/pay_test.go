package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
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

	orderPaid := &model.Order{
		OrderUUID:       orderUuid,
		UserUUID:        userUuid,
		PaymentMethod:   &paymentMethod,
		TransactionUUID: &transaction_uuid,
		OrderStatus:     model.OrderStatusPAID,
	}

	ctx := context.Background()
	s.orderRepo.On("GetOrder", ctx, orderUuid).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", ctx, orderUuid.String(), userUuid.String(), paymentMethod).Return(transaction_uuid, nil).Once()
	s.orderRepo.On("UpdateOrder", ctx, orderPaid).Return(nil)

	result, err := s.service.PayOrder(ctx, orderUuid, paymentMethod)

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

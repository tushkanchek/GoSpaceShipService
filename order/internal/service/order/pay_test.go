package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)




func (s *ServiceSuite) TestPayOrderSucces() {
	orderUuid := uuid.MustParse(gofakeit.UUID())

	userUuid := uuid.MustParse(gofakeit.UUID())

	paymentMethod := model.PaymentMethodCARD

	order := &model.Order{
		OrderUUID: orderUuid,
		UserUUID: userUuid,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}
	
	transaction_uuid := uuid.MustParse(gofakeit.UUID())

	orderPaid := &model.Order{
		OrderUUID: orderUuid,
		UserUUID: userUuid,
		PaymentMethod: &paymentMethod,
		TransactionUUID: &transaction_uuid,
		OrderStatus: model.OrderStatusPAID,
	}

	s.orderRepo.On("GetOrder", s.ctx, orderUuid).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUuid.String(), userUuid.String(), paymentMethod).Return(transaction_uuid, nil).Once()
	s.orderRepo.On("UpdateOrder", s.ctx, orderPaid).Return(nil)

	result, err := s.service.PayOrder(s.ctx, orderUuid, paymentMethod)

	s.NoError(err)
	s.Equal(transaction_uuid, result)
}




func (s *ServiceSuite) TestPayOrderEmptyOrderUuid() {
	orderUuid := uuid.Nil

	paymentMethod := model.PaymentMethodCARD


	result, err := s.service.PayOrder(s.ctx, orderUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())
	s.Empty(result)
}


package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v7"
)




func (s *ServiceSuite) TestPayOrderSucces() {
	orderUuid := gofakeit.UUID()

	userUuid := gofakeit.UUID()

	paymentMethod := model.PaymentMethodCARD

	order := &model.Order{
		OrderUUID: orderUuid,
		UserUUID: userUuid,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}
	
	transaction_uuid := gofakeit.UUID()

	orderPaid := &model.Order{
		OrderUUID: orderUuid,
		UserUUID: userUuid,
		PaymentMethod: &paymentMethod,
		TransactionUUID: &transaction_uuid,
		OrderStatus: model.OrderStatusPAID,
	}

	s.orderRepo.On("GetOrder", s.ctx, orderUuid).Return(order, nil).Once()
	s.paymentClient.On("PayOrder", s.ctx, orderUuid, userUuid, paymentMethod).Return(transaction_uuid, nil).Once()
	s.orderRepo.On("UpdateOrder", s.ctx, orderPaid).Return(nil).Once()

	result, err := s.service.PayOrder(s.ctx, orderUuid, paymentMethod)

	s.NoError(err)
	s.Equal(transaction_uuid, result)
}




func (s *ServiceSuite) TestPayOrderEmptyOrderUuid() {
	orderUuid := ""

	paymentMethod := model.PaymentMethodCARD


	result, err := s.service.PayOrder(s.ctx, orderUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())
	s.Empty(result)
}


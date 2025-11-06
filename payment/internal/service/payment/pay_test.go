package payment

import (
	"payment/internal/model"

	"github.com/brianvoe/gofakeit/v7"
)



func (s *ServiceSuite) TestPayOrderSucces() {
	orderUuid := gofakeit.UUID()

	userUuid := gofakeit.UUID()

	paymentMethod := "PAYMENT_METHOD_CARD"

	result, err := s.service.PayOrder(s.ctx, orderUuid, userUuid, paymentMethod)

	s.NoError(err)

	s.Require().NotNil(result)
}

func (s *ServiceSuite) TestPayOrderEmptyOrderUuid(){
	orderUuid := ""

	userUuid := gofakeit.UUID()

	paymentMethod := "PAYMENT_METHOD_CARD"

	result, err := s.service.PayOrder(s.ctx, orderUuid, userUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())

	s.Empty(result)
}


func (s *ServiceSuite) TestPayOrderEmptyUserUuid(){
	orderUuid := gofakeit.UUID()

	userUuid := ""

	paymentMethod := "PAYMENT_METHOD_CARD"

	result, err := s.service.PayOrder(s.ctx, orderUuid, userUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyUserUuid.Error())

	s.Empty(result)
}


func (s *ServiceSuite) TestPayOrderEmptyPaymentMethod(){
	orderUuid := gofakeit.UUID()

	userUuid := gofakeit.UUID()

	paymentMethod := ""

	result, err := s.service.PayOrder(s.ctx, orderUuid, userUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyPaymentMethod.Error())

	s.Empty(result)
}
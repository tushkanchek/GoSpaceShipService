package payment

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSucces() {
	orderUuid := gofakeit.UUID()

	userUuid := gofakeit.UUID()

	paymentMethod := "PAYMENT_METHOD_CARD"

	ctx := context.Background()
	result, err := s.service.PayOrder(ctx, orderUuid, userUuid, paymentMethod)

	s.NoError(err)

	s.Require().NotNil(result)
}

func (s *ServiceSuite) TestPayOrderEmptyOrderUuid() {
	orderUuid := ""

	userUuid := gofakeit.UUID()

	paymentMethod := "PAYMENT_METHOD_CARD"

	ctx := context.Background()
	result, err := s.service.PayOrder(ctx, orderUuid, userUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())

	s.Empty(result)
}

func (s *ServiceSuite) TestPayOrderEmptyUserUuid() {
	orderUuid := gofakeit.UUID()

	userUuid := ""

	paymentMethod := "PAYMENT_METHOD_CARD"

	ctx := context.Background()
	result, err := s.service.PayOrder(ctx, orderUuid, userUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyUserUuid.Error())

	s.Empty(result)
}

func (s *ServiceSuite) TestPayOrderEmptyPaymentMethod() {
	orderUuid := gofakeit.UUID()

	userUuid := gofakeit.UUID()

	paymentMethod := ""

	ctx := context.Background()
	result, err := s.service.PayOrder(ctx, orderUuid, userUuid, paymentMethod)

	s.EqualError(err, model.ErrEmptyPaymentMethod.Error())

	s.Empty(result)
}

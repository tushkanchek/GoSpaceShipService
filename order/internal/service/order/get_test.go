package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestGetOrderSucces() {
	order_uuid := uuid.MustParse(gofakeit.UUID())

	order := &model.Order{
		OrderUUID: order_uuid,
	}

	s.orderRepo.On("GetOrder", mock.Anything, order_uuid).Return(order, nil)

	ctx := context.Background()
	result, err := s.service.GetOrderByUUID(ctx, order_uuid)

	s.NoError(err)
	s.Equal(order, result)
}

func (s *ServiceSuite) TestGetOrderEmptyOrderUuid() {
	order_uuid := uuid.Nil

	ctx := context.Background()
	result, err := s.service.GetOrderByUUID(ctx, order_uuid)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())
	s.Nil(result)
}

func (s *ServiceSuite) TestGetOrderNotFound() {
	order_uuid := uuid.MustParse(gofakeit.UUID())

	s.orderRepo.On("GetOrder", mock.Anything, order_uuid).Return(nil, model.ErrOrderNotFound)

	ctx := context.Background()
	result, err := s.service.GetOrderByUUID(ctx, order_uuid)

	s.EqualError(err, model.ErrOrderNotFound.Error())
	s.Nil(result)
}

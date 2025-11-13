package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)


func (s *ServiceSuite) TestGetOrderSucces() {
	order_uuid := uuid.MustParse(gofakeit.UUID())

	order := &model.Order{
		OrderUUID: order_uuid,
	}

	s.orderRepo.On("GetOrder", s.ctx, order_uuid).Return(order, nil)

	result, err := s.service.GetOrderByUUID(s.ctx, order_uuid)

	s.NoError(err)
	s.Equal(order, result)
}


func (s *ServiceSuite) TestGetOrderEmptyOrderUuid() {
	order_uuid := uuid.Nil

	result, err := s.service.GetOrderByUUID(s.ctx, order_uuid)

	s.EqualError(err, model.ErrEmptyOrderUuid.Error())
	s.Nil(result)
}


func (s *ServiceSuite) TestGetOrderNotFound() {
	order_uuid := uuid.MustParse(gofakeit.UUID())

	s.orderRepo.On("GetOrder", s.ctx, order_uuid).Return(nil, model.ErrOrderNotFound)

	result, err := s.service.GetOrderByUUID(s.ctx, order_uuid)

	s.EqualError(err, model.ErrOrderNotFound.Error())
	s.Nil(result)
}



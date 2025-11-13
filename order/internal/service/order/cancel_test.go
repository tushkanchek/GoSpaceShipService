package order

import (
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)



func (s *ServiceSuite) TestCancelOrderSucces() {
	orderUuid := uuid.MustParse(gofakeit.UUID())

	order := &model.Order{
		OrderUUID: orderUuid,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepo.On("GetOrder", s.ctx, orderUuid).Return(order, nil).Once()
	s.orderRepo.On("UpdateOrder", s.ctx, 
		mock.MatchedBy(func(updatedOrder *model.Order) bool {
			return updatedOrder.OrderStatus == model.OrderStatusCANCELLED &&
				updatedOrder.OrderUUID == orderUuid
		}),
	).Return(nil).Once()
	
	err := s.service.CancelOrder(s.ctx, orderUuid)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderUuid := uuid.MustParse(gofakeit.UUID())

	s.orderRepo.On("GetOrder", s.ctx, orderUuid).Return(nil, model.ErrOrderNotFound).Once()
	
	err := s.service.CancelOrder(s.ctx, orderUuid)
	s.EqualError(err, model.ErrOrderNotFound.Error())
}



func (s *ServiceSuite) TestCancelOrderAlreadyPaid() {
	orderUuid := uuid.MustParse(gofakeit.UUID())
	order := &model.Order{
		OrderUUID: orderUuid,
		OrderStatus: model.OrderStatusPAID,
	}

	s.orderRepo.On("GetOrder", s.ctx, orderUuid).Return(order, nil).Once()

	err := s.service.CancelOrder(s.ctx, orderUuid)
	s.EqualError(err, model.ErrCancelOrderStatusPaid.Error())
}

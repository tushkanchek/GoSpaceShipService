package order

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order/internal/model"
)

func (s *ServiceSuite) TestCancelOrderSucces() {
	orderUuid := uuid.MustParse(gofakeit.UUID())

	order := &model.Order{
		OrderUUID:   orderUuid,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	ctx := context.Background()
	s.orderRepo.On("GetOrder", mock.Anything, orderUuid).Return(order, nil).Once()
	s.orderRepo.On("UpdateOrder", mock.Anything,
		mock.MatchedBy(func(updatedOrder *model.Order) bool {
			return updatedOrder.OrderStatus == model.OrderStatusCANCELLED &&
				updatedOrder.OrderUUID == orderUuid
		}),
	).Return(nil).Once()

	err := s.service.CancelOrder(ctx, orderUuid)
	s.NoError(err)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderUuid := uuid.MustParse(gofakeit.UUID())

	ctx := context.Background()
	s.orderRepo.On("GetOrder", mock.Anything, orderUuid).Return(nil, model.ErrOrderNotFound).Once()

	err := s.service.CancelOrder(ctx, orderUuid)
	s.EqualError(err, model.ErrOrderNotFound.Error())
}

func (s *ServiceSuite) TestCancelOrderAlreadyPaid() {
	orderUuid := uuid.MustParse(gofakeit.UUID())
	order := &model.Order{
		OrderUUID:   orderUuid,
		OrderStatus: model.OrderStatusPAID,
	}

	ctx := context.Background()
	s.orderRepo.On("GetOrder", mock.Anything, orderUuid).Return(order, nil).Once()

	err := s.service.CancelOrder(ctx, orderUuid)
	s.EqualError(err, model.ErrCancelOrderStatusPaid.Error())
}

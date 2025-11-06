package order

import (
	"math"
	"order/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateOrderSucces() {
	user_uuid := gofakeit.UUID()

	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()
	part_uuids := []string{uuid1, uuid2}
	price1 := gofakeit.Float64Range(0, 10000000)
	price2 := gofakeit.Float64Range(0, 10000000)

	part1 := &model.Part{
		Uuid:  uuid1,
		Price: price1,
	}

	part2 := &model.Part{
		Uuid:  uuid2,
		Price: price2,
	}

	parts := []*model.Part{part1, part2}

	filter := &model.PartsFilter{
		Uuids: part_uuids,
	}

	order := &model.Order{
		OrderUUID:   gofakeit.UUID(),
		UserUUID:    user_uuid,
		PartUuids:   part_uuids,
		TotalPrice:  price1 + price2,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(parts, nil).Once()

	s.orderRepo.On("CreateOrder", s.ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.UserUUID == user_uuid &&
			len(o.PartUuids) == len(part_uuids) &&
			price1+price2 == o.TotalPrice &&
			o.OrderStatus == model.OrderStatusPENDINGPAYMENT
	})).Return(nil).Once()

	result, err := s.service.CreateOrder(s.ctx, user_uuid, part_uuids)

	s.Nil(err)
	//s.Equal(order.OrderUUID, result.OrderUUID)
	s.Equal(order.UserUUID, result.UserUUID)
	s.Equal(order.OrderStatus, result.OrderStatus)
	s.Equal(order.TotalPrice, result.TotalPrice)

	s.inventoryClient.AssertExpectations(s.T())
	s.orderRepo.AssertExpectations(s.T())
}


func (s *ServiceSuite) TestCreateOrderAlreadyExists() {
	user_uuid := gofakeit.UUID()

	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()
	part_uuids := []string{uuid1, uuid2}
	price1 := gofakeit.Float64Range(0, 10000000)
	price2 := gofakeit.Float64Range(0, 10000000)

	part1 := &model.Part{
		Uuid:  uuid1,
		Price: price1,
	}

	part2 := &model.Part{
		Uuid:  uuid2,
		Price: price2,
	}

	parts := []*model.Part{part1, part2}

	filter := &model.PartsFilter{
		Uuids: part_uuids,
	}

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(parts, nil).Once()

	s.orderRepo.On("CreateOrder", s.ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.UserUUID == user_uuid &&
			len(o.PartUuids) == len(part_uuids) &&
			math.Abs(o.TotalPrice-(price1+price2)) < 0.0001 &&
			o.OrderStatus == model.OrderStatusPENDINGPAYMENT
	})).Return(model.ErrOrderAlreadyExists).Once()

	result, err := s.service.CreateOrder(s.ctx, user_uuid, part_uuids)

	s.Nil(result)
	s.EqualError(err, model.ErrOrderAlreadyExists.Error())
	

	s.inventoryClient.AssertExpectations(s.T())
	s.orderRepo.AssertExpectations(s.T())
}

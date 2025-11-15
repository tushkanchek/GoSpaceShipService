package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	model "inventory/internal/model"
)

// TODO: add more tests
func (s *ServiceSuite) TestListPartsSucces() {
	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()

	filter := &model.PartsFilter{
		Uuids:                 []string{uuid1, uuid2},
		Names:                 []string{"Jonathan Hoppe", "Albin Labadie"},
		Categories:            []model.Category{model.CategoryEngine, model.CategoryFuel},
		ManufacturerCountries: []string{"Estonia", "Canada"},
		Tags:                  []string{"engine", "premium"},
	}

	expectedParts := []*model.Part{
		{
			Uuid: uuid1,
		},
		{
			Uuid: uuid2,
		},
	}

	ctx := context.Background()
	s.inventoryRepo.On("ListParts", mock.Anything, filter).Return(expectedParts, nil)

	result, err := s.service.ListParts(ctx, filter)

	s.NoError(err)

	s.Equal(expectedParts, result)
}

func (s *ServiceSuite) TestListNamesFilterSucces() {
	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()
	filter := &model.PartsFilter{
		Names: []string{"Jonathan Hoppe", "Albin Labadie"},
	}

	expectedParts := []*model.Part{
		{
			Uuid: uuid1,
			Name: "Jonathan Hoppe",
		},
		{
			Uuid: uuid2,
			Name: "Albin Labadie",
		},
	}
	ctx := context.Background()
	s.inventoryRepo.On("ListParts", mock.Anything, filter).Return(expectedParts, nil)

	result, err := s.service.ListParts(ctx, filter)

	s.NoError(err)

	s.Equal(expectedParts, result)
}

func (s *ServiceSuite) TestListPartsNotFound() {
	filter := &model.PartsFilter{
		Uuids: []string{"uuid-unknown"},
	}

	ctx := context.Background()
	s.inventoryRepo.On("ListParts", mock.Anything, filter).Return(nil, model.ErrPartsNotFound)

	result, err := s.service.ListParts(ctx, filter)

	s.EqualError(err, model.ErrPartsNotFound.Error())

	s.Nil(result)
}

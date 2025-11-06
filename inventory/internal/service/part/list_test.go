package part

import (
	model "inventory/internal/model"

	repoConverter "inventory/internal/repository/converter"

	"github.com/brianvoe/gofakeit/v7"
)


//TODO: add more tests
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

	repoFilter := repoConverter.PartsFilterToRepoPartsFilter(filter)

	expectedParts := []*model.Part{
		{
			Uuid: uuid1,
		},
		{
			Uuid: uuid2,
		},
	}

	s.inventoryRepo.On("ListParts", s.ctx, repoFilter).Return(expectedParts, nil)
	
	result, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)

	s.Equal(expectedParts, result)

}


func (s *ServiceSuite) TestListNamesFilterSucces() {
	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()
	filter := &model.PartsFilter{
		Names:	[]string{"Jonathan Hoppe", "Albin Labadie"},
	}

	repoFilter := repoConverter.PartsFilterToRepoPartsFilter(filter)

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

	s.inventoryRepo.On("ListParts", s.ctx, repoFilter).Return(expectedParts, nil)
	
	result, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)

	s.Equal(expectedParts, result)

}


func (s *ServiceSuite) TestListPartsNotFound() {
	filter := &model.PartsFilter{
		Uuids:		[]string{"uuid-unknown"},
	}

	repoFilter := repoConverter.PartsFilterToRepoPartsFilter(filter)

	

	s.inventoryRepo.On("ListParts", s.ctx, repoFilter).Return(nil, model.ErrPartsNotFound)
	
	result, err := s.service.ListParts(s.ctx, filter)

	s.EqualError(err, model.ErrPartsNotFound.Error())

	s.Nil(result)
}
package part

import (
	"inventory/internal/model"

	"github.com/brianvoe/gofakeit/v7"
	
)

func (s *ServiceSuite) TestGetPartSucces() {
	partUuidStr := gofakeit.UUID()

	part := &model.Part{
		Uuid: partUuidStr,
	}

	s.inventoryRepo.On("GetPart", s.ctx, partUuidStr).Return(part, nil).Once()

	result, err := s.service.GetPart(s.ctx, partUuidStr)

	s.NoError(err)
	s.Equal(part, result)
}

func (s *ServiceSuite) TestGetPartEmptyUuid() {
	emptyUuid := ""

	result, err := s.service.GetPart(s.ctx, emptyUuid)

	s.EqualError(err, model.ErrPartUUIDIsEmpty.Error())
	s.Nil(result)
}

func (s *ServiceSuite) TestGetPartInvalidUuid() {
	invalidUuid := "apelsin-serega"

	result, err := s.service.GetPart(s.ctx, invalidUuid)

	s.EqualError(err, model.ErrUUIDIsNotValid.Error())
	s.Nil(result)
}

func (s *ServiceSuite) TestGetPartNotFoundPart(){
	partUuid := gofakeit.UUID()

	s.inventoryRepo.On("GetPart", s.ctx, partUuid).Return(&model.Part{}, model.ErrPartNotFound).Once()

	result, err := s.service.GetPart(s.ctx, partUuid)

	s.EqualError(err, model.ErrPartNotFound.Error())

	s.Nil(result)
	
}






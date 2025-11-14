package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSucces() {
	partUuidStr := gofakeit.UUID()

	part := &model.Part{
		Uuid: partUuidStr,
	}

	ctx := context.Background()
	s.inventoryRepo.On("GetPart", ctx, partUuidStr).Return(part, nil).Once()

	result, err := s.service.GetPart(ctx, partUuidStr)

	s.NoError(err)
	s.Equal(part, result)
}

func (s *ServiceSuite) TestGetPartEmptyUuid() {
	emptyUuid := ""

	ctx := context.Background()
	result, err := s.service.GetPart(ctx, emptyUuid)

	s.EqualError(err, model.ErrPartUUIDIsEmpty.Error())
	s.Nil(result)
}

func (s *ServiceSuite) TestGetPartInvalidUuid() {
	invalidUuid := "apelsin-serega"

	ctx := context.Background()
	result, err := s.service.GetPart(ctx, invalidUuid)

	s.EqualError(err, model.ErrUUIDIsNotValid.Error())
	s.Nil(result)
}

func (s *ServiceSuite) TestGetPartNotFoundPart() {
	partUuid := gofakeit.UUID()

	ctx := context.Background()
	s.inventoryRepo.On("GetPart", ctx, partUuid).Return(&model.Part{}, model.ErrPartNotFound).Once()

	result, err := s.service.GetPart(ctx, partUuid)

	s.EqualError(err, model.ErrPartNotFound.Error())

	s.Nil(result)
}

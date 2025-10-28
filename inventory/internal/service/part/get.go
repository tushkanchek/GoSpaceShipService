package part

import (
	"context"
	"inventory/internal/model"

	"github.com/google/uuid"
)

func (s *service) GetPart(ctx context.Context, partUuid string) (*model.Part, error){

	if _, err := uuid.Parse(partUuid);err!=nil{
		return nil, model.ErrUUIDIsNotValid
	}

	part, err := s.inventoryRepository.GetPart(ctx, partUuid)
	if err != nil{
		return nil, err
	}

	return part, err
}
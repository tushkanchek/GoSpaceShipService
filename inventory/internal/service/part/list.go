package part

import (
	"context"
	"inventory/internal/repository/converter"
	"inventory/internal/model"
)


func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error){
	return s.inventoryRepository.ListParts(ctx, converter.PartsFilterToRepoPartsFilter(filter))
}


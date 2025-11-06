package part

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	return s.inventoryRepository.ListParts(ctx, converter.PartsFilterToRepoPartsFilter(filter))
}

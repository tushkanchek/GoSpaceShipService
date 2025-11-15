package part

import (
	"context"

	"inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	reqCtxList, cancelList := context.WithTimeout(ctx, model.RequestTimeOutRead)
	defer cancelList()
	return s.inventoryRepository.ListParts(reqCtxList, filter)
}

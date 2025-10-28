package repository

import (
	"context"

	model "inventory/internal/model"
	repoModel "inventory/internal/repository/model"
)




type InventoryRepository interface{
	GetPart(ctx context.Context, partUuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *repoModel.PartsFilter) ([]*model.Part, error)
}
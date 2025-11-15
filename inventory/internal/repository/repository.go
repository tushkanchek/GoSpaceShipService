package repository

import (
	"context"

	model "inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, partUuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}

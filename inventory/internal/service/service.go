package service

import (
	"context"

	"inventory/internal/model"
)

//go:generate go run github.com/vektra/mockery/v3@latest --name=InventoryService
type InventoryService interface {
	GetPart(ctx context.Context, partUuid string) (*model.Part, error)
	ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}

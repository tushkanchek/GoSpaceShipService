package v1

import (
	"context"
	converter "inventory/internal/converter"
	inventoryV1 "shared/pkg/proto/inventory/v1"
)


func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error){
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err!=nil{
		return nil, err
	}
	return &inventoryV1.GetPartResponse{
		Part: converter.PartToApi(part),
	}, nil
}
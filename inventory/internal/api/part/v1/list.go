package v1

import (
	"context"
	"errors"
	converter "inventory/internal/converter"
	"inventory/internal/model"
	inventoryV1 "shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error){
	parts, err := a.inventoryService.ListParts(ctx, converter.PartsFilterToModel(req.Filter))
	if err!=nil{
		if errors.Is(err, model.ErrPartsNotFound){
			return nil, status.Errorf(codes.NotFound, "parts with this filter %s not found", req.GetFilter())
		}
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.PartsListToApi(parts),
	}, nil
}
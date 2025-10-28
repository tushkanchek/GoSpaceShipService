package v1

import (
	"context"
	"errors"
	converter "inventory/internal/converter"
	"inventory/internal/model"
	"log"
	inventoryV1 "shared/pkg/proto/inventory/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error){
	part, err := a.inventoryService.GetPart(ctx, req.GetUuid())
	if err!=nil{
		
		if errors.Is(err, model.ErrPartNotFound){
			log.Printf("GetPart: part with uuid %s not found\n", req.Uuid)
			return nil, status.Errorf(codes.NotFound, "part with uuid not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get part: %v", err)
	}
	return &inventoryV1.GetPartResponse{
		Part: converter.PartToApi(part),
	}, nil
}
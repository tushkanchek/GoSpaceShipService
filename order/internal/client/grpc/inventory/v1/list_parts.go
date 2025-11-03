package v1

import (
	"context"
	clientConverter "order/internal/client/converter"
	"order/internal/model"
	inventoryV1 "shared/pkg/proto/inventory/v1"
)


func (c *client) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error){
	//TODO: work with context

	resp, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err!=nil{
		return nil, err
	}


	parts := make([]*model.Part, 0, len(resp.Parts))
	for _, part := range resp.Parts{
		parts = append(parts, clientConverter.PartToModel(part))
	}
	return parts, nil
}
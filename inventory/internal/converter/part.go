package converter

import (
	model "inventory/internal/model"
	inventoryV1 "shared/pkg/proto/inventory/v1"

	"github.com/samber/lo"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// model "inventory/internal/model"
// repoModel "inventory/internal/repository/model"


func PartsFilterToModel(filter *inventoryV1.PartsFilter) *model.PartsFilter{
	categories := make([]model.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories{
		categories = append(categories, CategoryToModelCategory(category))
	}
	return &model.PartsFilter{
		Uuids: filter.Uuids,
		Names: filter.Names,
		Categories: categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags: filter.Tags,
	}
}


func CategoryToModelCategory(category inventoryV1.Category) model.Category{
	return model.Category(category)
}

func PartsListToApi(parts []*model.Part) []*inventoryV1.Part{
	apiParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts{
		apiParts = append(apiParts, PartToApi(part))
	}
	return apiParts
}


//be careful with type Value
func PartToApi(part *model.Part) *inventoryV1.Part{
	var dimensions *inventoryV1.Dimensions
	if part.Dimensions!=nil{
		dimensions = &inventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width: part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	var manafacturer *inventoryV1.Manufacturer
	if part.Manufacturer!=nil{
		manafacturer = &inventoryV1.Manufacturer{
			Name: part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	var createdAt *timestamppb.Timestamp
	if part.CreatedAt!=nil{
		createdAt = lo.ToPtr(*timestamppb.New(*part.CreatedAt))
	}

	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt!=nil{
		updatedAt = lo.ToPtr(*timestamppb.New(*part.UpdatedAt))
	}
	return &inventoryV1.Part{
		Uuid: part.Uuid,
		Name: part.Name,
		Description: part.Description,
		Price: part.Price,
		StockQuantity: part.StockQuantity,
		Category: inventoryV1.Category(part.Category),
		Dimensions: dimensions,
		Manufacturer: manafacturer,
		Tags: part.Tags,
		Metadata:metadataModelToProto(part.Metadata),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func metadataModelToProto(metadata map[string]model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		var protoValue *inventoryV1.Value
		switch {
		case value.StringValue != nil:
			protoValue = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{
					StringValue: *value.StringValue,
				},
			}
		case value.Int64Value != nil:
			protoValue = &inventoryV1.Value{
				Kind: &inventoryV1.Value_Int64Value{
					Int64Value: *value.Int64Value,
				},
			}
		case value.DoubleValue != nil:
			protoValue = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{
					DoubleValue: *value.DoubleValue,
				},
			}
		case value.BoolValue != nil:
			protoValue = &inventoryV1.Value{
				Kind: &inventoryV1.Value_BoolValue{
					BoolValue: *value.BoolValue,
				},
			}
		}
		if protoValue != nil {
			result[key] = protoValue
		}
	}
	return result
}

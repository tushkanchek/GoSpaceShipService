package converter

import (
	"log"

	"github.com/samber/lo"
	"order/internal/model"
	inventoryV1 "shared/pkg/proto/inventory/v1"
)

func PartToModel(part *inventoryV1.Part) *model.Part {
	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    DimensionToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata), // WARNING! model has any
		CreatedAt:     lo.ToPtr(part.CreatedAt.AsTime()),
		UpdatedAt:     lo.ToPtr(part.UpdatedAt.AsTime()),
	}
}

func DimensionToModel(dimension *inventoryV1.Dimensions) *model.Dimensions {
	return &model.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func ManufacturerToModel(manafacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	return &model.Manufacturer{
		Name:    manafacturer.Name,
		Country: manafacturer.Country,
		Website: manafacturer.Website,
	}
}

func MetadataToModel(protoMetadata map[string]*inventoryV1.Value) map[string]any {
	metadata := make(map[string]any, len(protoMetadata))
	for key, value := range protoMetadata {
		if value == nil || value.Kind == nil {
			continue
		}

		switch v := value.Kind.(type) {
		case *inventoryV1.Value_Int64Value:
			metadata[key] = v.Int64Value
		case *inventoryV1.Value_BoolValue:
			metadata[key] = v.BoolValue
		case *inventoryV1.Value_DoubleValue:
			metadata[key] = v.DoubleValue
		case *inventoryV1.Value_StringValue:
			metadata[key] = v.StringValue
		default:
			log.Printf("unkown metadata kind of key %q: %T", key, value.Kind)

		}
	}
	return metadata
}

func PartsFilterToProto(filter *model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}
	return &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

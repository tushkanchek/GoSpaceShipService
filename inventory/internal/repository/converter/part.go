package converter

import (
	"inventory/internal/model"
	repoModel "inventory/internal/repository/model"
)

func RepoPartToModel(part repoModel.Part) *model.Part {
	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      RepoCategoryToCategory(part.Category),
		Dimensions:    RepoDimensionToDimension(*part.Dimensions),
		Manufacturer:  RepoManufacturerToManufacturer(*part.Manufacturer),
		Tags:          part.Tags,
	}
}

func RepoCategoryToCategory(category repoModel.Category) model.Category {
	return model.Category(category)
}

func RepoDimensionToDimension(dimension repoModel.Dimensions) *model.Dimensions {
	if (dimension == repoModel.Dimensions{}) {
		return nil
	}

	return &model.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func RepoManufacturerToManufacturer(manafacturer repoModel.Manufacturer) *model.Manufacturer {
	if (manafacturer == repoModel.Manufacturer{}) {
		return nil
	}

	return &model.Manufacturer{
		Name:    manafacturer.Name,
		Country: manafacturer.Country,
		Website: manafacturer.Website,
	}
}

func PartsFilterToRepoPartsFilter(filter *model.PartsFilter) *repoModel.PartsFilter {
	if filter == nil {
		return nil
	}

	categories := make([]repoModel.Category, 0, len(filter.Categories))

	for _, category := range filter.Categories {
		categories = append(categories, CategoryToRepoCategory(model.Category(category)))
	}
	return &repoModel.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoryToRepoCategory(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

func DimensionToRepoDimension(dimension model.Dimensions) *repoModel.Dimensions {
	if (dimension == model.Dimensions{}) {
		return nil
	}

	return &repoModel.Dimensions{
		Length: dimension.Length,
		Width:  dimension.Width,
		Height: dimension.Height,
		Weight: dimension.Weight,
	}
}

func ManufacturerToRepoManufacturer(manafacturer model.Manufacturer) *repoModel.Manufacturer {
	if (manafacturer == model.Manufacturer{}) {
		return nil
	}

	return &repoModel.Manufacturer{
		Name:    manafacturer.Name,
		Country: manafacturer.Country,
		Website: manafacturer.Website,
	}
}

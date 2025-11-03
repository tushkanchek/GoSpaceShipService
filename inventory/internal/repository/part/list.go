package part

import (
	"context"

	model "inventory/internal/model"
	repoConverter "inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, Filter *repoModel.PartsFilter) ([]*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtredParts := FilterParts(r.data, Filter)

	if len(filtredParts) == 0 {
		return nil, model.ErrPartsNotFound
	}

	return filtredParts, nil
}

func FilterParts(parts map[string]*repoModel.Part, filter *repoModel.PartsFilter) []*model.Part {
	var PartsToReturn []*model.Part
	if filter == nil {
		PartsToReturn = make([]*model.Part, 0, len(parts))
		for i := range parts {
			PartsToReturn = append(PartsToReturn, repoConverter.RepoPartToModel(*parts[i]))
		}
		return PartsToReturn
	}

	for _, part := range parts {
		// Filter by UUID
		if len(filter.Uuids) > 0 && !contains(filter.Uuids, part.Uuid) {
			continue
		}
		if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
			continue
		}
		if len(filter.Categories) > 0 {
			found := false
			for _, v := range filter.Categories {
				if v == part.Category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if len(filter.Tags) > 0 && !anyMatch(part.Tags, filter.Tags) {
			continue
		}

		PartsToReturn = append(PartsToReturn, repoConverter.RepoPartToModel(*part))

	}
	return PartsToReturn
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// checks tags intersections empty or not
func anyMatch(tags, filter []string) bool {
	for _, tag := range tags {
		for _, f := range filter {
			if tag == f {
				return true
			}
		}
	}
	return false
}

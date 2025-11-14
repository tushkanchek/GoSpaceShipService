package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	model "inventory/internal/model"
	repoConverter "inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter *repoModel.PartsFilter) ([]*model.Part, error) {
	mongoFilter := bson.M{}
	if filter != nil {
		mongoFilter = buildMongoFilter(filter)
	}
	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor")
		}
	}()

	var parts []*model.Part
	for cursor.Next(ctx) {
		var p repoModel.Part
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		parts = append(parts, repoConverter.RepoPartToModel(p))
	}
	if len(parts) == 0 {
		log.Printf("found 0 parts according this filter")
		return nil, model.ErrPartsNotFound
	}

	return parts, nil
}

func buildMongoFilter(filter *repoModel.PartsFilter) bson.M {
	f := bson.M{}

	if len(filter.Uuids) > 0 {
		f["uuid"] = bson.M{"$in": filter.Uuids}
	}

	if len(filter.Names) > 0 {
		f["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		f["category"] = bson.M{"$in": filter.Categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		f["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		f["tags"] = bson.M{"$in": filter.Tags} // check any same tags
	}

	return f
}

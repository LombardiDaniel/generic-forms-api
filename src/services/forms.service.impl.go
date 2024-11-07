package services

import (
	"context"
	"log/slog"

	"github.com/LombardiDaniel/generic-data-collector-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FormsServiceMongoImpl struct {
	dataStoreCol *mongo.Collection
}

func NewFormsServiceMongoImpl(col *mongo.Collection) FormsService {
	return &FormsServiceMongoImpl{
		dataStoreCol: col,
	}
}

func (s *FormsServiceMongoImpl) InsertPayload(ctx context.Context, form models.Form) error {
	_, err := s.dataStoreCol.InsertOne(ctx, form)
	return err
}

func (s *FormsServiceMongoImpl) Get(ctx context.Context, id string) ([]models.Form, error) {
	query := bson.M{
		"id": id,
	}

	forms := []models.Form{}
	cur, err := s.dataStoreCol.Find(ctx, query)
	if err == mongo.ErrNilDocument {
		return forms, nil
	} else if err != nil {
		slog.Error(err.Error())
		return forms, err
	}
	err = cur.All(ctx, &forms)
	if err != nil {
		slog.Error(err.Error())
		return forms, err
	}

	return forms, err
}

func (s *FormsServiceMongoImpl) GetN(ctx context.Context, id string, n uint32) ([]models.Form, error) {
	query := bson.M{
		"id": id,
	}

	forms := []models.Form{}

	opts := options.Find().SetLimit(int64(n))

	cur, err := s.dataStoreCol.Find(ctx, query, opts)
	if err == mongo.ErrNilDocument {
		return forms, nil
	} else if err != nil {
		slog.Error(err.Error())
		return forms, err
	}
	err = cur.All(ctx, &forms)
	if err != nil {
		slog.Error(err.Error())
		return forms, err
	}

	return forms, err
}

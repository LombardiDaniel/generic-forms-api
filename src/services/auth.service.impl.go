package services

import (
	"context"

	"github.com/LombardiDaniel/generic-data-collector-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthServiceImpl struct {
	tokensCol *mongo.Collection
}

func NewAuthServiceImpl(tokensCol *mongo.Collection) AuthService {
	return &AuthServiceImpl{
		tokensCol: tokensCol,
	}
}

func (s *AuthServiceImpl) Authenticate(ctx context.Context, key string) error {
	var token models.Token

	query := bson.M{
		"token": key,
	}

	err := s.tokensCol.FindOne(ctx, query).Decode(&token)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceImpl) CreateToken(ctx context.Context, token models.Token) error {
	_, err := s.tokensCol.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	return err
}

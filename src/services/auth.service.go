package services

import (
	"context"

	"github.com/LombardiDaniel/generic-data-collector-api/models"
)

type AuthService interface {
	Authenticate(ctx context.Context, key string) error
	CreateToken(ctx context.Context, token models.Token) error
}

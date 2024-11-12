package services

import (
	"context"

	"github.com/LombardiDaniel/generic-data-collector-api/models"
)

type FormService interface {
	InsertPayload(ctx context.Context, formPayload models.Form) error
	Get(ctx context.Context, id string) ([]models.Form, error)
	GetN(ctx context.Context, id string, n uint32) ([]models.Form, error)
}

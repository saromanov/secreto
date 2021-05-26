package storage

import (
	"context"
	"errors"

	"github.com/saromanov/secreto/internal/models"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	CreateSecret(ctx context.Context, secret models.Secret) error
	GetSecret(ctx context.Context, key string) (models.Secret, error)
}

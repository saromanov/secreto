package badger

import (
	"context"
	"os"
	"testing"

	"github.com/saromanov/secreto/internal/models"
	"github.com/saromanov/secreto/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecret(t *testing.T) {
	ctx := context.Background()
	s, err := New(storage.Config{
		Path:          "/tmp/test1",
		EncryptionKey: "test",
	})
	assert.NoError(t, err)
	assert.NoError(t, s.CreateSecret(ctx, models.Secret{
		Key:   "key",
		Value: "value",
	}))
	assert.Error(t, s.CreateSecret(ctx, models.Secret{
		Key:   "",
		Value: "value",
	}))
	assert.NoError(t, os.RemoveAll("/tmp/test1"))
}

func TestGetSecret(t *testing.T) {
	ctx := context.Background()
	s, err := New(storage.Config{
		Path:          "/tmp/test1",
		EncryptionKey: "test",
	})
	assert.NoError(t, err)
	s.CreateSecret(ctx, models.Secret{
		Key:   "key",
		Value: "value",
	})
	assert.NoError(t, os.RemoveAll("/tmp/test1"))

	res, err := s.GetSecret(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, "value", res.Value)

	res, err = s.GetSecret(ctx, "keya")
	assert.Error(t, err)
	assert.Equal(t, "", res.Value)
	assert.NoError(t, os.RemoveAll("/tmp/test1"))
}
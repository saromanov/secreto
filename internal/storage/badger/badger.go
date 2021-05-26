package badger

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	"github.com/saromanov/secreto/internal/models"
	"github.com/saromanov/secreto/internal/storage"
)

type store struct {
	db *badger.DB
}

// New provides initialization for badger
func New(cfg storage.Config) (storage.Storage, error) {
	opts := badger.DefaultOptions("/tmp/badger")
	opts.EncryptionKey = []byte("test")
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to open badger")
	}
	return &store{
		db: db,
	}, nil
}

func (d *store) CreateSecret(ctx context.Context, secret models.Secret) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(secret.Key), []byte(secret.Value))
		return err
	})
	if err != nil {
		return errors.Wrap(err, "unable to set data")
	}
	return nil
}

func (d *store) GetSecret(ctx context.Context, key string) (models.Secret, error) {
	var value []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		if err := item.Value(func(val []byte) error {
			value = val
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return models.Secret{}, err
	}
	return models.Secret{
		Key:   key,
		Value: string(value),
	}, nil
}

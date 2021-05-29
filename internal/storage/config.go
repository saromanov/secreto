package storage

// Config defines path for store data
type Config struct {
	Path          string `env:"BADGER_PATH"`
	EncryptionKey string `env:"BADGER_ENCRYPTIONKEY"`
}

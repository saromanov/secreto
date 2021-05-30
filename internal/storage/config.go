package storage

// Config defines path for store data
type Config struct {
	Path          string `env:"BADGER_PATH,default=/tmp/badger"`
	EncryptionKey string `env:"BADGER_ENCRYPTIONKEY,default=test"`
}

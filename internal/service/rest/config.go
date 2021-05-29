package rest

// Config defines configuration for rest
type Config struct {
	Port int `env:"REST_PORT,default=8089"`
}
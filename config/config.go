package config

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	conf Config
	once sync.Once
)

type Config struct {
	Port            int           `envconfig:"PORT" default:"8080"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"5s"`
	ProviderTimeout time.Duration `envconfig:"PROVIDER_TIMEOUT" default:"200ms"`

	RedisHost string `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort int    `envconfig:"REDIS_PORT" default:"6379"`
}

func Get() Config {
	once.Do(func() {
		envconfig.MustProcess("", &conf)
	})
	return conf
}

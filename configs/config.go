package configs

import (
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port            int           `envconfig:"PORT" default:"8080"`
	ProviderTimeout time.Duration `envconfig:"PROVIDER_TIMEOUT" default:"200ms"`

	RedisHost string `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort int    `envconfig:"REDIS_PORT" default:"6379"`
}

func Get() Config {
	conf := Config{}
	once := sync.Once{}

	once.Do(func() {
		envconfig.MustProcess("", &conf)
	})

	return conf
}

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
	ServerPort            int           `envconfig:"SERVER_PORT" default:"8080"`
	ServerShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
	ServerReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"5s"`
	ServerWriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	ServerIdleTimeout  time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"60s"`

	ProviderTimeout      time.Duration `envconfig:"PROVIDER_TIMEOUT" default:"500ms"`
	ProviderMaxRetries   int           `envconfig:"PROVIDER_MAX_RETRIES" default:"3"`
	ProviderRetryDelay   time.Duration `envconfig:"PROVIDER_RETRY_DELAY" default:"100ms"`
	ProviderRateLimitRPS float64       `envconfig:"PROVIDER_RATE_LIMIT_RPS" default:"10"`

	DefaultCacheTTL time.Duration `envconfig:"DEFAULT_CACHE_TTL" default:"300s"`

	RedisHost string `envconfig:"REDIS_HOST" default:"redis"`
	RedisPort int    `envconfig:"REDIS_PORT" default:"6379"`
}

func Get() Config {
	once.Do(func() {
		envconfig.MustProcess("", &conf)
	})
	return conf
}

package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type retryEnvConfig struct {
	MaxRetries        int           `env:"RETRY_MAX_ATTEMPTS,required"`
	InitialDelay      time.Duration `env:"RETRY_INITIAL_DELAY,required"`
	MaxDelay          time.Duration `env:"RETRY_MAX_DELAY,required"`
	BackoffMultiplier float64       `env:"RETRY_BACKOFF_MULTIPLIER,required"`
}

type retryConfig struct {
	raw retryEnvConfig
}

func NewRetryConfig() (*retryConfig, error) {
	var raw retryEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &retryConfig{
		raw: raw,
	}, nil
}

func (cfg *retryConfig) MaxRetries() int {
	return cfg.raw.MaxRetries
}

func (cfg *retryConfig) InitialDelay() time.Duration {
	return cfg.raw.InitialDelay
}

func (cfg *retryConfig) MaxDelay() time.Duration {
	return cfg.raw.MaxDelay
}

func (cfg *retryConfig) BackoffMultiplier() float64 {
	return cfg.raw.BackoffMultiplier
}

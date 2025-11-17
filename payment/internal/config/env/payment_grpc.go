package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type PaymentGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type paymentGRPCConfig struct {
	raw PaymentGRPCEnvConfig
}

func NewPaymentGRPCConfig() (*paymentGRPCConfig, error) {
	var raw PaymentGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &paymentGRPCConfig{
		raw: raw,
	}, nil
}

func (cfg *paymentGRPCConfig) Adress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

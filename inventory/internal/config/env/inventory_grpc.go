package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type InventoryGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw InventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw InventoryGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{
		raw: raw,
	}, nil
}

func (cfg *inventoryGRPCConfig) Adress() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

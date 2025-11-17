package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderHTTPConfig interface {
	Adress() string
	ReadTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
	MigrationDir() string
}

type PaymentGRPCConfig interface {
	Adress() string
}

type InventoryGRPCConfig interface {
	Adress() string
}

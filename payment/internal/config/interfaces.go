package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PaymentGRPCConfig interface {
	Adress() string
}

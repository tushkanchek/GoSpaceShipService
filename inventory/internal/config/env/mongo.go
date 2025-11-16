package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)


type MongoEnvConfig struct{
	ImageName			string		`env:"MONGO_IMAGE_NAME,required"`
	ExternalPort		string		`env:"EXTERNAL_MONGO_PORT,required"`
	Host				string		`env:"MONGO_HOST,required"`
	Port				string		`env:"MONGO_PORT,required"`
	Database			string		`env:"MONGO_DATABASE,required"`
	AuthDB				string		`env:"MONGO_AUTH_DB,required"`
	Username			string		`env:"MONGO_INITDB_ROOT_USERNAME,required"`
	Password			string		`env:"MONGO_INITDB_ROOT_PASSWORD,required"`
}


type mongoConfig struct{
	raw MongoEnvConfig
}

func NewMongoConfig() (*mongoConfig, error) {
	var raw MongoEnvConfig
	if err:=env.Parse(&raw);err!=nil {
		return nil, err
	}

	return &mongoConfig{
		raw: raw,
	}, nil
}


func (cfg *mongoConfig) URI() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=%s",
		cfg.raw.Username,
		cfg.raw.Password,
		cfg.raw.Host,
		cfg.raw.Port,
		cfg.raw.Database,
		cfg.raw.AuthDB,
	)
}

func (cfg *mongoConfig) DatabaseName() string {
	return cfg.raw.Database
}


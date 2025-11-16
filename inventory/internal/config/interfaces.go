package config


type LoggerConfig interface{
	Level() string
	AsJson() bool
}


type InventoryGRPCConfig interface{
	Adress() string
}

type MongoConfig interface{
	URI() string
	DatabaseName() string
}
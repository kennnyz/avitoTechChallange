package config

import (
	"os"
)

type Config struct {
	DB         *DBConfig
	ServerAddr string `json:"server_addr"`
	SwaggerURL string `json:"swagger_url"`
}

type DBConfig struct {
	Dsn string
}

func ReadConfig() (*Config, error) {
	var cfg Config

	db := &DBConfig{}
	cfg.DB = db
	cfg.DB.Dsn = os.Getenv("DB_DSN")
	cfg.ServerAddr = os.Getenv("SERVER_ADDR")
	cfg.SwaggerURL = os.Getenv("SWAGGER_URL")

	return &cfg, nil
}

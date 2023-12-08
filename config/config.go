package config

import (
	"fmt"
	env "github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	Nats        Nats
	DataBase    DB
	Cache       Cache
}

type Nats struct {
	URL string `env:"NATS_URL" envDefault:"127.0.0.1:4222"`
}

type Cache struct {
	Addr string `env:"CACHE_ADRR" envDefault:"127.0.0.1:6379"`
}

type DB struct {
	DbType         string `env:"DB_TYPE" envDefault:"postgres://"`
	DbUser         string `env:"POSTGRES_USER" envDefault:"userL0:"`
	DbUserPassword string `env:"POSTGRES_PASSWORD" envDefault:"123456"`
	DbPort         string `env:"DB_PORT" envDefault:"@0.0.0.0:5432/"`
	DbName         string `env:"POSTGRES_DB" envDefault:"postgres"`
}

func GetConfig() *Config {
	var appConfig Config

	if err := env.Parse(&appConfig); err != nil {
		fmt.Printf("read configuration error: %s\n", err)
		os.Exit(1)
	}

	return &appConfig
}

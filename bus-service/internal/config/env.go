package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environment struct {
	Port          int    `env:"PORT" envDefault:"8080"`
	PostgresURL   string `env:"POSTGRES_URL"`
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func LoadEnvironment() (*Environment, error) {
	_ = godotenv.Load()

	cfg := &Environment{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

package config

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	RPC          string         `env:"RPC" env-required:"true"`
	TokenAddress common.Address `env:"TOKEN_ADDRESS" env-required:"true"`
	DB           string         `env:"DB" env-required:"true"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env: %w", err)
	}
	return cfg, nil
}

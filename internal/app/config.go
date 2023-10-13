package app

import (
	"github.com/kelseyhightower/envconfig"
	"testtask_frankrg/internal/server"
)

type AppConfig struct {
	ServerConfig server.ServerConfig `envconfig:"SERVER"`
}

func newConfig() (*AppConfig, error) {
	cfg := new(AppConfig)
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

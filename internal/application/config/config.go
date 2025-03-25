package config

import (
	"flag"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/caarlos0/env/v6"
)

type AppConfig struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
}

type ConfigValidater interface {
	Validate(opts *AppConfig) error
}

func LoadConfig(validater ConfigValidater) (*AppConfig, error) {
	opts, err := ParseFlags(validater)
	if err != nil {
		return nil, err
	}
	err = ParseEnvConfig(opts)
	if err != nil {
		return nil, err
	}

	log.Zap.Infof("Config loaded: %+v\n", opts)
	return opts, nil
}

func ParseFlags(validater ConfigValidater) (*AppConfig, error) {
	opts := new(AppConfig)

	flag.StringVar(&opts.RunAddress, "a", ":8080", "listen address")
	flag.StringVar(&opts.AccrualSystemAddress, "d", "http://localhost:8081", "accrual system URI")
	flag.StringVar(&opts.DatabaseURI, "r", "", "database connection string")

	flag.Parse()

	err := validater.Validate(opts)

	if err != nil {
		return nil, fmt.Errorf("config validating failed: %w", err)
	}

	return opts, nil
}

func ParseEnvConfig(opts *AppConfig) error {
	err := env.Parse(opts)

	if err != nil {
		return fmt.Errorf("failed env parsing: %w", err)
	}

	return nil
}

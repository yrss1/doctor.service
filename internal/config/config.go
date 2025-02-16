package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configs struct {
	APP      AppConfig
	POSTGRES StoreConfig
}

type AppConfig struct {
	Mode    string
	Port    string
	Path    string
	Timeout time.Duration
}

type StoreConfig struct {
	DSN string
}

func New() (*Configs, error) {
	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	if err := godotenv.Load(filepath.Join(root, ".env")); err != nil {
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	var cfg Configs

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	return &cfg, nil
}

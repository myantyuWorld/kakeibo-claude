package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port     string
	Origins  []string
	Postgres Postgres
}

func Load() (*Config, error) {
	port, err := getMustEnv("PORT")
	if err != nil {
		return nil, err
	}

	pg, err := loadPostgres()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:     port,
		Origins:  strings.Split(getEnvOrDefault("ORIGINS", ""), ","),
		Postgres: pg,
	}, nil
}

func getMustEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("environment variable %q is required", key)
	}
	return v, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

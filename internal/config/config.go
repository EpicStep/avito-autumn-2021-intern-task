package config

import (
	"fmt"
	"os"
)

// Config struct.
type Config struct {
	Port      string
	PgURL     string
	EAPIToken string
}

// New config.
func New() (*Config, error) {
	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}

	pgURL, err := getEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}

	eAPIToken, err := getEnv("EXCHANGERATESAPI_TOKEN")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:      port,
		PgURL:     pgURL,
		EAPIToken: eAPIToken,
	}, nil
}

func getEnv(key string) (string, error) {
	value, isFounded := os.LookupEnv(key)
	if isFounded {
		return value, nil
	}

	return "", fmt.Errorf("env variable %s not presented", key)
}

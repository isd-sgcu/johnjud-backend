package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port        int
	Env         string
	MaxFileSize int64
}

type ServiceConfig struct {
	Auth    string
	Backend string
	File    string
}

type Config struct {
	App     AppConfig
	Service ServiceConfig
}

func LoadConfig() (*Config, error) {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	port, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)
	if err != nil {
		return nil, err
	}
	maxFileSize, err := strconv.ParseInt(os.Getenv("APP_MAX_FILE_SIZE"), 10, 64)
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{
		Port:        int(port),
		Env:         os.Getenv("APP_ENV"),
		MaxFileSize: maxFileSize,
	}

	serviceConfig := ServiceConfig{
		Auth:    os.Getenv("SERVICE_AUTH"),
		Backend: os.Getenv("SERVICE_BACKEND"),
		File:    os.Getenv("SERVICE_FILE"),
	}

	return &Config{
		App:     appConfig,
		Service: serviceConfig,
	}, nil

}

func (ac *AppConfig) IsDevelopment() bool {
	return ac.Env == "development"
}

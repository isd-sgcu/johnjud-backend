package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type App struct {
	Port        int
	Env         string
	MaxFileSize int64
}

type Service struct {
	Auth    string
	Backend string
	File    string
}

type Database struct {
	Url string
}

type Redis struct {
	Host     string
	Port     int
	Password string
}

type Jwt struct {
	Secret          string
	ExpiresIn       int
	RefreshTokenTTL int
	Issuer          string
	ResetTokenTTL   int
}

type Auth struct {
	ClientURL string
}

type Sendgrid struct {
	ApiKey  string
	Name    string
	Address string
}

type Config struct {
	App      App
	Service  Service
	Database Database
	Redis    Redis
	Jwt      Jwt
	Auth     Auth
	Sendgrid Sendgrid
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
	app := App{
		Port:        int(port),
		Env:         os.Getenv("APP_ENV"),
		MaxFileSize: maxFileSize,
	}

	service := Service{
		Auth:    os.Getenv("SERVICE_AUTH"),
		Backend: os.Getenv("SERVICE_BACKEND"),
		File:    os.Getenv("SERVICE_FILE"),
	}

	database := Database{
		Url: os.Getenv("DATABASE_URL"),
	}

	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return nil, err
	}
	redis := Redis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     redisPort,
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	jwtExpiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return nil, err
	}
	jwtRefreshTokenTTL, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}
	jwtResetTokenTTL, err := strconv.Atoi(os.Getenv("JWT_RESET_TOKEN_TTL"))
	if err != nil {
		return nil, err
	}
	jwt := Jwt{
		Secret:          os.Getenv("JWT_SECRET"),
		ExpiresIn:       jwtExpiresIn,
		RefreshTokenTTL: jwtRefreshTokenTTL,
		Issuer:          os.Getenv("JWT_ISSUER"),
		ResetTokenTTL:   jwtResetTokenTTL,
	}

	auth := Auth{
		ClientURL: os.Getenv("AUTH_CLIENT_URL"),
	}

	sendgrid := Sendgrid{
		ApiKey:  os.Getenv("SENDGRID_API_KEY"),
		Name:    os.Getenv("SENDGRID_NAME"),
		Address: os.Getenv("SENDGRID_ADDRESS"),
	}

	return &Config{
		App:      app,
		Service:  service,
		Database: database,
		Redis:    redis,
		Jwt:      jwt,
		Auth:     auth,
		Sendgrid: sendgrid,
	}, nil

}

func (ac *App) IsDevelopment() bool {
	return ac.Env == "development"
}

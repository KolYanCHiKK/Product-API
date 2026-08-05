package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db   *DbConfig
	Auth *AuthConfig
}

type DbConfig struct {
	DSN string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file, using default config")
	}

	return &Config{
		Db: &DbConfig{
			DSN: os.Getenv("DSN"),
		},
		Auth: &AuthConfig{
			Secret: os.Getenv("TOKEN"),
		},
	}
}

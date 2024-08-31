package config

import (
	"ctserver/dotenv"
	"os"
)

type Config struct {
	Port        string
	DatabaseURI string
}

func New() (*Config, error) {
	err := dotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURI: os.Getenv("DATABASE_URI"),
	}, nil
}

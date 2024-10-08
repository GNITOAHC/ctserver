package config

import (
	"ctserver/dotenv"
	"os"
)

type Config struct {
	Port             string
	DatabaseURI      string
	SMTPFrom         string
	SMTPPass         string
	SMTPHost         string
	SMTPPort         string
	OTPSecret        string
	AuthDBURI        string
	AuthDBName       string
	AuthDBCollection string
	JWTSecret        string
	BaseURL          string

	TestMail string // Testing Environment
}

func New(envpaths ...string) (*Config, error) {
	// Load environment variables from .env file
	err := dotenv.Load(envpaths...)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:             os.Getenv("PORT"),
		DatabaseURI:      os.Getenv("DATABASE_URI"),
		SMTPFrom:         os.Getenv("SMTP_FROM"),
		SMTPPass:         os.Getenv("SMTP_PASS"),
		SMTPHost:         os.Getenv("SMTP_HOST"),
		SMTPPort:         os.Getenv("SMTP_PORT"),
		OTPSecret:        os.Getenv("OTP_SECRET"),
		AuthDBURI:        os.Getenv("AUTH_DB_URI"),
		AuthDBName:       os.Getenv("AUTH_DB_NAME"),
		AuthDBCollection: os.Getenv("AUTH_DB_COLLECTION"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		BaseURL:          os.Getenv("BASE_URL"),

		TestMail: os.Getenv("TEST_MAIL"), // Testing Environment
	}, nil
}

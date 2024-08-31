package config

import (
	"ctserver/dotenv"
	"os"
)

type Config struct {
	Port        string
	DatabaseURI string
	SMTPFrom    string
	SMTPPass    string
	SMTPHost    string
	SMTPPort    string
	OTPSecret   string

	TestMail string // Testing Environment
}

func New(envpaths ...string) (*Config, error) {
	// Load environment variables from .env file
	err := dotenv.Load(envpaths...)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		DatabaseURI: os.Getenv("DATABASE_URI"),
		SMTPFrom:    os.Getenv("SMTP_FROM"),
		SMTPPass:    os.Getenv("SMTP_PASS"),
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		OTPSecret:   os.Getenv("OTP_SECRET"),

		TestMail: os.Getenv("TEST_MAIL"), // Testing Environment
	}, nil
}

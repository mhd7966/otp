package main

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App struct {
		Debug      bool   `env:"DEBUG" envDefault:"true"`
		Env        string `env:"ENV" envDefault:"development"`
		InstanceID string `env:"INSTANCE_ID" envDefault:"otp-service"`
	} `envPrefix:"APP_"`

	Database struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     int    `env:"PORT" envDefault:"5432"`
		User     string `env:"USER" envDefault:"postgres"`
		Password string `env:"PASSWORD" envDefault:"password"`
		DBName   string `env:"DB_NAME" envDefault:"otp_db"`
		SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
	} `envPrefix:"DB_"`

	Redis struct {
		Host     string `env:"HOST" envDefault:"localhost"`
		Port     int    `env:"PORT" envDefault:"6379"`
		Password string `env:"PASSWORD" envDefault:""`
		DB       int    `env:"DB" envDefault:"0"`
	} `envPrefix:"REDIS_"`

	Log struct {
		SentryDSN string `env:"SENTRY_DSN" envDefault:""`
	} `envPrefix:"LOG_"`
}

// Configs loads configuration from environment variables
func Configs(envFile string) (*Config, error) {
	// Load .env file if it exists
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return nil, err
		}
	}

	config := &Config{}

	// Set default values
	config.App.Debug = true
	config.App.Env = "development"
	config.App.InstanceID = "otp-service"

	config.Database.Host = "localhost"
	config.Database.Port = 5432
	config.Database.User = "postgres"
	config.Database.Password = "password"
	config.Database.DBName = "otp_db"
	config.Database.SSLMode = "disable"

	config.Redis.Host = "localhost"
	config.Redis.Port = 6379
	config.Redis.Password = ""
	config.Redis.DB = 0

	// Override with environment variables if they exist
	if env := os.Getenv("APP_DEBUG"); env != "" {
		config.App.Debug = env == "true"
	}
	if env := os.Getenv("APP_ENV"); env != "" {
		config.App.Env = env
	}
	if env := os.Getenv("DB_HOST"); env != "" {
		config.Database.Host = env
	}
	if env := os.Getenv("DB_PORT"); env != "" {
		// You might want to parse this as int
	}
	if env := os.Getenv("DB_USER"); env != "" {
		config.Database.User = env
	}
	if env := os.Getenv("DB_PASSWORD"); env != "" {
		config.Database.Password = env
	}
	if env := os.Getenv("DB_NAME"); env != "" {
		config.Database.DBName = env
	}
	if env := os.Getenv("REDIS_HOST"); env != "" {
		config.Redis.Host = env
	}
	if env := os.Getenv("REDIS_PORT"); env != "" {
		// You might want to parse this as int
	}

	return config, nil
}

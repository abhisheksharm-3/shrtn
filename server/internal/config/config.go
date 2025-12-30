// Package config provides application configuration loading and validation.
package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Environment        string
	ServerAddress      string
	AppwriteEndpoint   string
	AppwriteProjectID  string
	AppwriteAPIKey     string
	AppwriteCollection string
	AppwriteDatabase   string
	CORSOrigins        []string
	APIKey             string
	RateLimitPerMinute int
	RateLimitBurst     int
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Environment:        getEnv("ENVIRONMENT", "development"),
		ServerAddress:      getEnv("SERVER_ADDRESS", ":8080"),
		AppwriteEndpoint:   getEnv("APPWRITE_ENDPOINT", ""),
		AppwriteProjectID:  getEnv("APPWRITE_PROJECT_ID", ""),
		AppwriteAPIKey:     getEnv("APPWRITE_API_KEY", ""),
		AppwriteCollection: getEnv("APPWRITE_COLLECTION_ID", ""),
		AppwriteDatabase:   getEnv("APPWRITE_DATABASE_ID", ""),
		CORSOrigins:        parseCORSOrigins(getEnv("CORS_ORIGINS", "http://localhost:5173")),
		APIKey:             getEnv("API_KEY", ""),
		RateLimitPerMinute: getEnvInt("RATE_LIMIT_PER_MINUTE", 60),
		RateLimitBurst:     getEnvInt("RATE_LIMIT_BURST", 10),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.AppwriteProjectID == "" {
		return errors.New("APPWRITE_PROJECT_ID is required")
	}
	if c.AppwriteAPIKey == "" {
		return errors.New("APPWRITE_API_KEY is required")
	}
	if c.AppwriteCollection == "" {
		return errors.New("APPWRITE_COLLECTION_ID is required")
	}
	if c.AppwriteDatabase == "" {
		return errors.New("APPWRITE_DATABASE_ID is required")
	}
	return nil
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func parseCORSOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

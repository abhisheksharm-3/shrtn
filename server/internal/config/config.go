package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress      string
	AppwriteEndpoint   string
	AppwriteProjectID  string
	AppwriteAPIKey     string
	AppwriteCollection string
	AppwriteDatabase   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	// Get environment variables with default values
	return &Config{
		ServerAddress:      getEnv("SERVER_ADDRESS", ":8080"),
		AppwriteEndpoint:   getEnv("APPWRITE_ENDPOINT", ""),
		AppwriteProjectID:  getEnv("APPWRITE_PROJECT_ID", ""),
		AppwriteAPIKey:     getEnv("APPWRITE_API_KEY", ""),
		AppwriteCollection: getEnv("APPWRITE_COLLECTION_ID", ""),
		AppwriteDatabase:   getEnv("APPWRITE_DATABASE_ID", ""),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if defaultValue == "" {
		fmt.Printf("Warning: Required environment variable %s is not set\n", key)
	}

	return defaultValue
}

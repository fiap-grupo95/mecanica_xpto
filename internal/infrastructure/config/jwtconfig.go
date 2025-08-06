package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetSecretKey() string {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET not set in the .env file or environment variables")
	}
	return secretKey
}

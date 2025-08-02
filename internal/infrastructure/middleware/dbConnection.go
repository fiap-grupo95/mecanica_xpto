package middleware

import (
	"log"
	"mecanica_xpto/internal/infrastructure/config"

	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	db, err := config.NewDBFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err.Error())
	}
	return db
}

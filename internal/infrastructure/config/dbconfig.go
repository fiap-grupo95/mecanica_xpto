package config

import (
	"fmt"
	"mecanica_xpto/internal/domain/model/dto"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBFromEnv() (*gorm.DB, error) {
	// Carrega variáveis do .env
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&dto.PartsSupplyDTO{},
		&dto.ServiceDTO{},
		&dto.VehicleDTO{},
	)
	return db, nil
}

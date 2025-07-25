package middleware

import (
	"database/sql"
	"log"
	"mecanica_xpto/internal/config"
)

func ConnectDatabase() *sql.DB {
	db, err := config.NewDBFromEnv()
	defer db.Close()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err.Error())
	}
	return db
}

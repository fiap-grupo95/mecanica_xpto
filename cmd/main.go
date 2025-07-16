package main

import (
	"mecanica_xpto/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.NewDBFromEnv()
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	defer db.Close()

	// Configura o roteador Gin
	r := gin.Default()

	// Exemplo de rota GET
	r.GET("/ping", func(c *gin.Context) {
		// Aqui você pode usar o db
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Inicia o servidor na porta 8080
	r.Run(":8080")
}

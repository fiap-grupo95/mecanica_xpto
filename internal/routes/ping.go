package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func addPingRoutes(rg *gin.RouterGroup, db *sql.DB) {
	ping := rg.Group(PathHealthCheck)

	ping.GET("/", func(c *gin.Context) {
		// Aqui você pode usar o db
		//if err := db.Ping(); err != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"message": "DB error"})
		//	return
		//}
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

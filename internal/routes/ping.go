package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func addPingRoutes(rg *gin.RouterGroup) {
	ping := rg.Group(PathHealthCheck)

	ping.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
}

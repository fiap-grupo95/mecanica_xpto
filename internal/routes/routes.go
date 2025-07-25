package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"mecanica_xpto/internal/middleware"
	"strconv"
)

var router = gin.Default()

const port = 8080

// Run will start the server
func Run() {
	setMiddlewares()

	getRoutes()

	err := router.Run(":" + strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Failed to startup the application: %v", err.Error())
	}
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file,
// so this one won't be so messy
func getRoutes() {
	v1 := router.Group("/v1")
	addPingRoutes(v1)

	// other routes can be added here
	//addUserRoutes(v1, db)
}

func setMiddlewares() {
	// Set trusted proxies
	middleware.SetTrustedProxies(router)
	middleware.ConnectDatabase()

	// Set CORS middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))
}

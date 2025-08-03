package routes

import (
	"log"
	_ "mecanica_xpto/docs" // This will be auto-generated
	middleware2 "mecanica_xpto/internal/infrastructure/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

const PORT = 8080

// Run will start the server
func Run() {
	setMiddlewares()

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	getRoutes()

	err := router.Run(":" + strconv.Itoa(PORT))
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
	addUserRoutes(v1)
}

// setMiddlewares will configure our middleware
func setMiddlewares() {
	// Set trusted proxies
	middleware2.SetTrustedProxies(router)
	middleware2.ConnectDatabase()

	// Set CORS middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))
}

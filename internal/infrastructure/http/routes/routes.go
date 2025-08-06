package routes

import (
	"log"
	_ "mecanica_xpto/docs" // This will be auto-generated
	"mecanica_xpto/internal/infrastructure/config"
	database "mecanica_xpto/internal/infrastructure/databse"
	"mecanica_xpto/internal/infrastructure/http/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var router = gin.Default()

const PORT = 8080

// Run will start the server
func Run() {
	db := setMiddlewares()

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	getRoutes(db)

	err := router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to startup the application: %v", err.Error())
	}
}

// getRoutes will create our routes of our entire application
// this way every group of routes can be defined in their own file,
// so this one won't be so messy
func getRoutes(db *gorm.DB) {
	v1 := router.Group("/v1")
	addPingRoutes(v1)
	addUserRoutes(v1)
	addPartsSupplyRoutes(v1, db)
}

// setMiddlewares will configure our middleware
func setMiddlewares() *gorm.DB {

	secretKey := config.GetSecretKey()

	router.Use(middleware.JWTAuthMiddleware(secretKey))
	// Set trusted proxies
	middleware.SetTrustedProxies(router)
	db := database.ConnectDatabase()

	// Set CORS middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))

	return db
}

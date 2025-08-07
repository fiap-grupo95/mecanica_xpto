package routes

import (
	"log"
	_ "mecanica_xpto/docs" // This will be auto-generated
	repository "mecanica_xpto/internal/domain/repository/parts_supply"
	use_case "mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/internal/infrastructure/config"
	database "mecanica_xpto/internal/infrastructure/databse"
	"mecanica_xpto/internal/infrastructure/http"
	"mecanica_xpto/internal/infrastructure/http/middleware"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

const PORT = 8080

// Run will start the server
func Run() {
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

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	partsSupplyUseCase := use_case.NewPartsSupplyUseCase(repository.NewPartsSupplyRepository(db))
	partsSupplyHandler := http.NewPartsSupplyHandler(partsSupplyUseCase)

	v1 := router.Group("/v1")
	addPingRoutes(v1)
	addUserRoutes(v1)
	addPartsSupplyRoutes(v1, partsSupplyHandler)

	err := router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to startup the application: %v", err.Error())
	}
}

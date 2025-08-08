package routes

import (
	"log"
	"mecanica_xpto/internal/domain/repository/parts_supply"
	"mecanica_xpto/internal/domain/repository/vehicles"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/internal/infrastructure/database"
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
	//secretKey := config.GetSecretKey()
	//
	//router.Use(middleware.JWTAuthMiddleware(secretKey))
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

	partsSupplyUseCase := usecase.NewPartsSupplyUseCase(parts_supply.NewPartsSupplyRepository(db))
	partsSupplyHandler := http.NewPartsSupplyHandler(partsSupplyUseCase)

	vehiclesRepository := vehicles.NewVehicleRepository(db)
	vehiclesUseCase := usecase.NewVehicleService(vehiclesRepository)
	vehicleHandler := http.NewVehicleHandler(vehiclesUseCase)

	v1 := router.Group("/v1")
	addPingRoutes(v1)
	addPartsSupplyRoutes(v1, partsSupplyHandler)
	addVehicleRoutes(v1, vehicleHandler)

	err := router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to startup the application: %v", err.Error())
	}
}

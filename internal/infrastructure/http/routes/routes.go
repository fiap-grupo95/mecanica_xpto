package routes

import (
	"log"
	"mecanica_xpto/internal/domain/repository"
	memory "mecanica_xpto/internal/domain/repository/user-example/repository"
	"mecanica_xpto/internal/domain/service"
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

	userRepo := memory.NewMemoryRepository()
	userHandler := http.NewUserHandler(userRepo)

	vehiclesRepository := repository.NewVehicleRepository(db)
	vehiclesService := service.NewVehicleService(vehiclesRepository)
	vehicleHandler := http.NewVehicleHandler(vehiclesService)

	v1 := router.Group("/v1")
	addPingRoutes(v1)
	addUserRoutes(v1, userHandler)
	addVehicleRoutes(v1, vehicleHandler)

	err := router.Run(":" + strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("Failed to startup the application: %v", err.Error())
	}
}

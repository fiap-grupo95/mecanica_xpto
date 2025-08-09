package routes

import (
	"log"
	_ "mecanica_xpto/docs" // This will be auto-generated
	"mecanica_xpto/internal/domain/repository/parts_supply"
	"mecanica_xpto/internal/domain/repository/service"
	"mecanica_xpto/internal/domain/repository/vehicles"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/internal/infrastructure/database"
	"mecanica_xpto/internal/infrastructure/http"
	"mecanica_xpto/internal/infrastructure/http/handlers"
	"mecanica_xpto/internal/infrastructure/http/middleware"
	"mecanica_xpto/pkg/utils"
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

func getRoutes() {
	// Config JWT
	jwtCfg := utils.LoadJWTConfig()
	jwtService := utils.NewJWTService(jwtCfg)

	db := database.ConnectDatabase()

	// Handler de autenticação
	authHandler := handlers.NewAuthHandler(
		usecase.NewAuthUseCase(jwtService),
	)

	// Rotas públicas
	v1 := router.Group("/v1")
	v1.POST("/login", authHandler.Login)

	partsSupplyUseCase := usecase.NewPartsSupplyUseCase(parts_supply.NewPartsSupplyRepository(db))
	partsSupplyHandler := http.NewPartsSupplyHandler(partsSupplyUseCase)

	serviceUseCase := usecase.NewServiceUseCase(service.NewServiceRepository(db))
	serviceHandler := http.NewServiceHandler(serviceUseCase)

	vehiclesRepository := vehicles.NewVehicleRepository(db)
	vehiclesUseCase := usecase.NewVehicleService(vehiclesRepository)
	vehicleHandler := http.NewVehicleHandler(vehiclesUseCase)

	// Rotas protegidas
	authGroup := v1.Group("/")
	authGroup.Use(middleware.AuthMiddleware(jwtService))
	addPingRoutes(authGroup)
	addPartsSupplyRoutes(authGroup, partsSupplyHandler)
	addVehicleRoutes(authGroup, vehicleHandler)
	addServiceRoutes(authGroup, serviceHandler)
}

func setMiddlewares() {

	middleware.SetTrustedProxies(router)

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Recovered from panic: %v", recovered)
		c.AbortWithStatus(500)
	}))
}

package routes

import (
	"mecanica_xpto/internal/domain/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mecanica_xpto/internal/domain/repository"
	"mecanica_xpto/internal/infrastructure/http"
)

func addVehicleRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	vehicles := rg.Group(PathVehicles)
	vehiclesRepository := repository.NewVehicleRepository(db)
	vehiclesService := service.NewVehicleService(vehiclesRepository)
	vehicleHandler := http.NewVehicleHandler(vehiclesService)
	{
		vehicles.GET("/", vehicleHandler.GetVehicles)
	}
}

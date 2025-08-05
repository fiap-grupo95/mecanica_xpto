package routes

import (
	"github.com/gin-gonic/gin"
	"mecanica_xpto/internal/infrastructure/http"
)

func addVehicleRoutes(rg *gin.RouterGroup) {
	vehicles := rg.Group(PathVehicles)
	vehicleHandler := http.NewVehicleHandler()
	{

		vehicles.GET("/:id", vehicleHandler.GetVehicle)
	}
}

package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addVehicleRoutes(rg *gin.RouterGroup, vehicleHandler *http.VehicleHandler) {
	vehicles := rg.Group(PathVehicles)
	{
		vehicles.GET("/", vehicleHandler.GetVehicles)
		vehicles.GET("/customer/:customerID", vehicleHandler.GetVehiclesByCustomerID)
		vehicles.GET("/:id", vehicleHandler.GetVehicleByID)
		vehicles.GET("/:plate", vehicleHandler.GetVehicleByPlate)
		vehicles.POST("/", vehicleHandler.CreateVehicle)
		vehicles.PATCH("/:id", vehicleHandler.UpdateVehicle)
		vehicles.DELETE("/:id", vehicleHandler.DeleteVehicle)
	}
}

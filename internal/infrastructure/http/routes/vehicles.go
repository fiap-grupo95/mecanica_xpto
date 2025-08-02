package routes

import "github.com/gin-gonic/gin"

func addVehicleRoutes(rg *gin.RouterGroup) {
	vehicles := rg.Group(PathVehicles)
	{

		vehicles.GET("/:id", vehicleHandler.GetVehicle)
	}
}

package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addPartsSupplyRoutes(rg *gin.RouterGroup, partsSupplyHandler *http.PartsSupplyHandler) {

	partsSupply := rg.Group(PathPartsSupply)
	{
		partsSupply.GET("/:id", partsSupplyHandler.GetPartsSupplyByID)
		partsSupply.GET("/", partsSupplyHandler.ListPartsSupplies)
		partsSupply.POST("/", partsSupplyHandler.CreatePartsSupply)
		partsSupply.PUT("/:id", partsSupplyHandler.UpdatePartsSupply)
		partsSupply.DELETE("/:id", partsSupplyHandler.DeletePartsSupply)
	}
}

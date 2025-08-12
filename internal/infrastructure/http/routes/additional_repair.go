package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addAdditionalRepairRoutes(rg *gin.RouterGroup, additionalRepair *http.AdditionalRepairHandler) {
	serviceOrdersRoutes := rg.Group(PathAdditionalRepair)
	{
		serviceOrdersRoutes.POST("", additionalRepair.CreateSOAdditionalRepair)
		serviceOrdersRoutes.GET("/:id", additionalRepair.GetAdditionalRepair)
		serviceOrdersRoutes.PATCH("/:id/add", additionalRepair.AddPartSupplyAndService)
		serviceOrdersRoutes.PATCH("/:id/remove", additionalRepair.RemovePartSupplyAndService)
	}
}

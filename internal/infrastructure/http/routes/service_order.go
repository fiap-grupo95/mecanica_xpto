package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addServiceOrderRoutes(rg *gin.RouterGroup, serviceOrderHandler *http.ServiceOrderHandler) {
	serviceOrdersRoutes := rg.Group(PathServiceOrders)
	{
		serviceOrdersRoutes.GET("/:id", serviceOrderHandler.GetServiceOrder)
		serviceOrdersRoutes.POST("", serviceOrderHandler.CreateServiceOrder)
		serviceOrdersRoutes.PATCH("/:id/diagnosis", serviceOrderHandler.UpdateServiceOrderDiagnosis)
		serviceOrdersRoutes.PATCH("/:id/estimate", serviceOrderHandler.UpdateServiceOrderEstimate)
		serviceOrdersRoutes.PATCH("/:id/execution", serviceOrderHandler.UpdateServiceOrderExecution)
		serviceOrdersRoutes.PATCH("/:id/delivery", serviceOrderHandler.UpdateServiceOrderDelivery)
		serviceOrdersRoutes.GET("/", serviceOrderHandler.ListServiceOrders)
	}
}

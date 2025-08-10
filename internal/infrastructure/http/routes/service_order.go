package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addServiceOrderRoutes(rg *gin.RouterGroup, serviceOrderHandler *http.ServiceOrderHandler) {
	serviceOrdersRoutes := rg.Group(PathServiceOrders)
	{
		//serviceOrdersRoutes.GET("/:id", serviceOrderHandler.GetServiceOrder)
		serviceOrdersRoutes.POST("", serviceOrderHandler.CreateServiceOrder)
		serviceOrdersRoutes.PATCH("/:id", serviceOrderHandler.UpdateServiceOrder)
		//serviceOrdersRoutes.GET("/", serviceOrderHandler.ListServiceOrders)
	}
}

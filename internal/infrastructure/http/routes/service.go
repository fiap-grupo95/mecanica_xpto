package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addServiceRoutes(rg *gin.RouterGroup, serviceHandler *http.ServiceHandler) {

	service := rg.Group(PathService)
	{
		service.GET("/:id", serviceHandler.GetServiceByID)
		service.GET("/", serviceHandler.ListServices)
		service.POST("/", serviceHandler.CreateService)
		service.PUT("/:id", serviceHandler.UpdateService)
		service.DELETE("/:id", serviceHandler.DeleteService)
	}
}

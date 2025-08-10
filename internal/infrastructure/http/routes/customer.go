package routes

import (
	"github.com/gin-gonic/gin"
	"mecanica_xpto/internal/infrastructure/http"
)

func addCustomerRoutes(rg *gin.RouterGroup, customerHandler *http.CustomerHandler) {
	customersRoutes := rg.Group(PathCustomers)
	{
		customersRoutes.GET("/full/:id", customerHandler.GetFullCustomer)
		customersRoutes.GET("/:document", customerHandler.GetCustomer)
		customersRoutes.POST("", customerHandler.CreateCustomer)
		customersRoutes.PATCH("/:id", customerHandler.UpdateCustomer)
		customersRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
		customersRoutes.GET("/", customerHandler.ListCustomer)
	}
}

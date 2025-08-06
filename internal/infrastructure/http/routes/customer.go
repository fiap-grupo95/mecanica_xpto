package routes

import (
	"mecanica_xpto/internal/domain/repository/customers"
	database "mecanica_xpto/internal/infrastructure/databse"
	"mecanica_xpto/internal/infrastructure/http/handler/customers"

	"github.com/gin-gonic/gin"
)

func addCustomerRoutes(rg *gin.RouterGroup) {
	db := database.ConnectDatabase()

	customerRepo := customers.NewCustomerRepository(db)
	customerHandler := http.NewCustomerHandler(customerRepo)

	customersRoutes := rg.Group(PathCustomers)
	{
		customersRoutes.GET("/:id", customerHandler.GetCustomer)
	}
}

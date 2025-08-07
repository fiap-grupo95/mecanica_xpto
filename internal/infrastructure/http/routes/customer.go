package routes

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/repository/customers"
	use_cases "mecanica_xpto/internal/domain/use_cases/customer"
	database "mecanica_xpto/internal/infrastructure/databse"
	"mecanica_xpto/internal/infrastructure/http/handler/customers"

	"github.com/gin-gonic/gin"
)

func addCustomerRoutes(rg *gin.RouterGroup) {
	db := database.ConnectDatabase()
	err := db.AutoMigrate(&dto.CustomerDTO{})
	err = db.AutoMigrate(&dto.UserDTO{})
	err = db.AutoMigrate(&dto.VehicleDTO{})
	if err != nil {
		return
	}

	userRepo := repository.NewUserRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	customerUseCase := use_cases.NewCustomerUseCase(customerRepo, userRepo)
	customerHandler := http.NewCustomerHandler(customerUseCase)

	customersRoutes := rg.Group(PathCustomers)
	{
		customersRoutes.GET("/full/:id", customerHandler.GetFullCustomer)
		customersRoutes.GET("/:document", customerHandler.GetCustomer)
		customersRoutes.POST("", customerHandler.CreateCustomer)
		customersRoutes.PATCH("/:id", customerHandler.UpdateCustomer)
		customersRoutes.DELETE("/:id", customerHandler.DeleteCustomer)
		customersRoutes.GET("/all", customerHandler.ListCustomer)
	}
}

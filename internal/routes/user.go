package routes

import (
	"mecanica_xpto/internal/domain/user/repository"
	"mecanica_xpto/internal/handler"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup) {
	userRepo := repository.NewMemoryRepository()
	userHandler := handler.NewUserHandler(userRepo)

	users := rg.Group(PathUsers)
	{
		users.GET("/:id", userHandler.GetUser)
	}
}

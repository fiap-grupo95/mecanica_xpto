package routes

import (
	"mecanica_xpto/internal/domain/user-example/repository"
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup) {
	userRepo := repository.NewMemoryRepository()
	userHandler := http.NewUserHandler(userRepo)

	users := rg.Group(PathUsers)
	{
		users.GET("/:id", userHandler.GetUser)
	}
}

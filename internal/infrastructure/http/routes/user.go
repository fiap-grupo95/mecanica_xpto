package routes

import (
	"mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
)

func addUserRoutes(rg *gin.RouterGroup, userHandler *http.UserHandler) {
	users := rg.Group(PathUsers)
	{
		users.GET("/:id", userHandler.GetUser)
	}
}

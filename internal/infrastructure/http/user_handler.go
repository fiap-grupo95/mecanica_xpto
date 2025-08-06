package http

import (
	"mecanica_xpto/internal/domain/repository/user-example"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo user_example.Repository
}

// NewUserHandler creates a new user-example http
func NewUserHandler(repo user_example.Repository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUser godoc
// @Summary Get user-example by ID
// @Description Retrieve a user-example by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user_example.User
// @Failure 404 {object} map[string]string "error":"user-example not found"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	foundUser, err := h.repo.GetByID(id)
	if err == user_example.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "user-example not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

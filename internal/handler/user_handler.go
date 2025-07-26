package handler

import (
	"github.com/gin-gonic/gin"
	"mecanica_xpto/internal/domain/user"
	"net/http"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo user.Repository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo user.Repository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.User
// @Failure 404 {object} map[string]string "error":"user not found"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	foundUser, err := h.repo.GetByID(id)
	if err == user.ErrUserNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, foundUser)
}

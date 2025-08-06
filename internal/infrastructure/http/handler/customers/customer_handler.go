package http

import (
	"mecanica_xpto/internal/domain/repository/customers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CustomerHandler handles HTTP requests for users
type CustomerHandler struct {
	repo customers.ICustomerRepository
}

// NewCustomerHandler creates a new customer http handler
func NewCustomerHandler(repo customers.ICustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

// GetCustomer godoc
// @Summary Get customer by ID
// @Description Retrieve a customer by their ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} customer.Customer
// @Failure 404 {object} map[string]string "error":"customer not found"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	foundCustomer, err := h.repo.GetByID(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, foundCustomer)
}

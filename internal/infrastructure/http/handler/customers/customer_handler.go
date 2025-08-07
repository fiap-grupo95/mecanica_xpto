package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	use_cases "mecanica_xpto/internal/domain/use_cases/customer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CustomerHandler handles HTTP requests for users
type CustomerHandler struct {
	ucCustomer use_cases.ICustomerUseCase
}

// NewCustomerHandler creates a new customer http handler
func NewCustomerHandler(us use_cases.ICustomerUseCase) *CustomerHandler {
	return &CustomerHandler{ucCustomer: us}
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
// @Router /customers/{document} [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	doc := c.Param("document")

	foundCustomer, err := h.ucCustomer.GetByDocument(doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, foundCustomer)
}

func (h *CustomerHandler) GetFullCustomer(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	foundCustomer, err := h.ucCustomer.GetById(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, foundCustomer)
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Creates a new customer record
// @Tags customers
// @Accept json
// @Produce json
// @Param vehicle body entities.Customer true "Customer information"
// @Success 201 {object} entities.Customer
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "error message"
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var body entities.Customer

	if err := c.BindJSON(&body); err != nil {
		return
	}

	err := h.ucCustomer.CreateCustomer(&body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

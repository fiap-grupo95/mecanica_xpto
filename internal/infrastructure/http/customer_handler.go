package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	use_cases "mecanica_xpto/internal/domain/usecase"
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

// GetFullCustomer godoc
// @Summary Get full customer by ID
// @Description Retrieve a full customer record by their numeric ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} entities.Customer
// @Failure 400 {object} map[string]string "error":"invalid customer id"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/id/{id} [get]
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

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update an existing customer record by ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body entities.Customer true "Customer information"
// @Success 200 {object} nil
// @Failure 400 {object} map[string]string "error":"invalid customer id or input"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	var body entities.Customer
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err = h.ucCustomer.UpdateCustomer(uint(idUint), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Delete a customer record by ID
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string "error":"invalid customer id"
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	err = h.ucCustomer.DeleteCustomer(uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListCustomer godoc
// @Summary List all customers
// @Description Retrieve a list of all customers
// @Tags Customers
// @Accept json
// @Produce json
// @Success 200 {array} entities.Customer
// @Failure 500 {object} map[string]string "error":"internal server error"
// @Router /customers [get]
func (h *CustomerHandler) ListCustomer(c *gin.Context) {
	customers, err := h.ucCustomer.ListCustomer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

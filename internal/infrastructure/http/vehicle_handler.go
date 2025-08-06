package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VehicleHandler handles HTTP requests for vehicle operations
// @title Vehicle API
type VehicleHandler struct {
	service service.VehicleServiceInterface
}

// NewVehicleHandler creates a new vehicle handler instance
func NewVehicleHandler(service service.VehicleServiceInterface) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

// GetVehicles godoc
// @Summary Get all vehicles
// @Description Retrieves a list of all vehicles
// @Tags vehicles
// @Accept json
// @Produce json
// @Success 200 {array} entities.Vehicle
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles [get]
func (v VehicleHandler) GetVehicles(c *gin.Context) {
	vehicles, err := v.service.GetAllVehicles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

// GetVehiclesByCustomerID godoc
// @Summary Get vehicles by customer ID
// @Description Retrieves all vehicles belonging to a specific customer
// @Tags vehicles
// @Accept json
// @Produce json
// @Param customerID path int true "Customer ID"
// @Success 200 {array} entities.Vehicle
// @Failure 400 {object} map[string]string "Invalid customer ID"
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles/customer/{customerID} [get]
func (v VehicleHandler) GetVehiclesByCustomerID(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("customerID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	vehicles, err := v.service.GetVehiclesByCustomerID(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

// CreateVehicle godoc
// @Summary Create a new vehicle
// @Description Creates a new vehicle record
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vehicle body entities.Vehicle true "Vehicle information"
// @Success 201 {object} entities.Vehicle
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles [post]
func (v VehicleHandler) CreateVehicle(c *gin.Context) {
	var vehicle entities.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := v.service.CreateVehicle(vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// UpdateVehicle godoc
// @Summary Update a vehicle
// @Description Updates an existing vehicle's information
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Param vehicle body entities.Vehicle true "Vehicle information"
// @Success 200 {object} entities.Vehicle
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles/{id} [put]
func (v VehicleHandler) UpdateVehicle(c *gin.Context) {
	var vehicle entities.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}
	vehicle.ID = uint(id)

	result, err := v.service.UpdateVehicle(vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Deletes a vehicle by its ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Invalid vehicle ID"
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles/{id} [delete]
func (v VehicleHandler) DeleteVehicle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	if err := v.service.DeleteVehicle(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

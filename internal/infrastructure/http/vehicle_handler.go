package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// VehicleHandler handles HTTP requests for vehicle operations
// @title Vehicle API
type VehicleHandler struct {
	service usecase.VehicleServiceInterface
}

// NewVehicleHandler creates a new vehicle handler instance
func NewVehicleHandler(service usecase.VehicleServiceInterface) *VehicleHandler {
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

// GetVehicleByCustomerID godoc
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
func (v VehicleHandler) GetVehicleByCustomerID(c *gin.Context) {
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

// GetVehicleByID godoc
// @Summary Get vehicles by ID
// @Description Retrieves a vehicle belonging to a specific ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param ID path int true "ID"
// @Summary Get vehicle by ID
// @Description Retrieves a vehicle by its ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param ID path int true "ID"
// @Success 200 {object} entities.Vehicle
// @Failure 400 {object} map[string]string "Invalid vehicle ID"
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles/{id} [get]
func (v VehicleHandler) GetVehicleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	vehicles, err := v.service.GetVehicleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// GetVehicleByPlate godoc
// @Summary Get vehicle by plate
// @Description Retrieves a vehicle belonging to a specific plate
// @Tags vehicles
// @Accept json
// @Produce json
// @Param string path string true "plate"
// @Success 200 {array} entities.Vehicle
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 500 {object} map[string]string "error message"
// @Router /vehicles/plate/{plate} [get]
func (v VehicleHandler) GetVehicleByPlate(c *gin.Context) {
	plate := c.Param("plate")
	if plate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plate cannot be empty"})
		return
	}
	vehicles, err := v.service.GetVehicleByPlate(plate)
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
// @Param vehicle body map[string]string true "Vehicle information"
// @Success 201 {object} string "Vehicle created successfully"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "error message"
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
// @Summary Update a vehicle partially
// @Description Updates specific fields of an existing vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Param vehicle body entities.Vehicle true "Vehicle fields to update"
// @Success 200 {object} string "Vehicle updated successfully"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "error message"
// @Router /vehicles/{id} [patch]
func (v VehicleHandler) UpdateVehicle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := v.service.UpdateVehiclePartial(uint(id), updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result})
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

package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	service service.VehicleServiceInterface
}

func NewVehicleHandler(service service.VehicleServiceInterface) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

func (v VehicleHandler) GetVehicles(c *gin.Context) {
	vehicles, err := v.service.GetAllVehicles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

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

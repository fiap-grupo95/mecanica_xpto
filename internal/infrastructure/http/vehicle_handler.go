package http

import (
	"mecanica_xpto/internal/domain/service"

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
		c.JSON(500, gin.H{})
	}
	c.JSON(200, vehicles)
}

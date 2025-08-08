package http

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	usecase "mecanica_xpto/internal/domain/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	usecase usecase.IServiceUseCase
}

func NewServiceHandler(usecase usecase.IServiceUseCase) *ServiceHandler {
	return &ServiceHandler{usecase: usecase}
}

func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	service, err := h.usecase.GetServiceByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, service)
}

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service entities.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	createdService, err := h.usecase.CreateService(c.Request.Context(), &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdService)
}

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var service entities.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	service.ID = uint(id)
	err = h.usecase.UpdateService(c.Request.Context(), &service)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
}

func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	err = h.usecase.DeleteService(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, usecase.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *ServiceHandler) ListServices(c *gin.Context) {
	services, err := h.usecase.ListServices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

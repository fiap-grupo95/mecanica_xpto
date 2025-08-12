package http

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	usecase "mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// Domain errors
	errInvalidServiceID    = pkg.NewDomainErrorSimple("INVALID_SERVICE_ID", "Invalid service ID", http.StatusBadRequest)
	errInvalidServiceInput = pkg.NewDomainErrorSimple("INVALID_INPUT", "Invalid input data", http.StatusBadRequest)
)

// ServiceHandler handles HTTP requests for service operations
// @title Service API
// @version 1.0
// @description API for managing services in the workshop management system
type ServiceHandler struct {
	usecase usecase.IServiceUseCase
}

func NewServiceHandler(usecase usecase.IServiceUseCase) *ServiceHandler {
	return &ServiceHandler{usecase: usecase}
}

func mapServiceError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrServiceNotFound):
		return pkg.NewDomainErrorSimple("SERVICE_NOT_FOUND", "Service not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid service ID", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrServiceAlreadyExists):
		return pkg.NewDomainErrorSimple("SERVICE_EXISTS", "Service already exists", http.StatusConflict)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetServiceByID godoc
// @Summary Get service by ID
// @Description Retrieve a service by its ID
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} entities.Service
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services/{id} [get]
func (h *ServiceHandler) GetServiceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(errInvalidServiceID.HTTPStatus, errInvalidServiceID.ToHTTPError())
		return
	}
	service, err := h.usecase.GetServiceByID(c.Request.Context(), uint(id))
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, service)
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service record
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param service body entities.Service true "Service Information"
// @Success 201 {object} entities.Service
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services [post]
func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service entities.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(errInvalidServiceInput.HTTPStatus, errInvalidServiceInput.ToHTTPError())
		return
	}
	createdService, err := h.usecase.CreateService(c.Request.Context(), &service)
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusCreated, createdService)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service record
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param service body entities.Service true "Service Information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services/{id} [put]
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(errInvalidServiceID.HTTPStatus, errInvalidServiceID.ToHTTPError())
		return
	}
	var service entities.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(errInvalidServiceInput.HTTPStatus, errInvalidServiceInput.ToHTTPError())
		return
	}
	service.ID = uint(id)
	err = h.usecase.UpdateService(c.Request.Context(), &service)
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully"})
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete an existing service record
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(errInvalidServiceID.HTTPStatus, errInvalidServiceID.ToHTTPError())
		return
	}
	err = h.usecase.DeleteService(c.Request.Context(), uint(id))
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// ListServices godoc
// @Summary List all services
// @Description Get a list of all services
// @Tags Services
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} entities.Service
// @Failure 500 {object} pkg.AppError
// @Router /services [get]
func (h *ServiceHandler) ListServices(c *gin.Context) {
	services, err := h.usecase.ListServices(c.Request.Context())
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, services)
}

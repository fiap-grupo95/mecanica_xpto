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
	errInvalidServiceID = pkg.NewDomainErrorSimple("INVALID_SERVICE_ID", "Invalid service ID", http.StatusBadRequest)
)

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

func (h *ServiceHandler) CreateService(c *gin.Context) {
	var service entities.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
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

func (h *ServiceHandler) UpdateService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(errInvalidServiceID.HTTPStatus, errInvalidServiceID.ToHTTPError())
		return
	}
	var service entities.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
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

func (h *ServiceHandler) ListServices(c *gin.Context) {
	services, err := h.usecase.ListServices(c.Request.Context())
	if err != nil {
		appErr := mapServiceError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, services)
}

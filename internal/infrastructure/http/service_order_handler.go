package http

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// operation flow
const (
	DIAGNOSIS = "diagnosis"
	ESTIMATE  = "estimate"
	EXECUTION = "execution"
	DELIVERY  = "delivery"
)

// ServiceOrderHandler handles HTTP requests related to service orders.
// It provides methods to create, update, retrieve, and list service orders.
// @title Service Order API
// @version 1.0
// @description API for managing service orders in the workshop management system
type ServiceOrderHandler struct {
	serviceOrderUseCase usecase.IServiceOrderUseCase
}

func NewServiceOrderHandler(useCase usecase.IServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		serviceOrderUseCase: useCase,
	}
}

// CreateServiceOrder godoc
// @Summary Create a new service order
// @Description Create a new service order record
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param order body entities.ServiceOrder true "Service Order Information"
// @Success 201 {object} entities.ServiceOrder
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders [post]
func (h *ServiceOrderHandler) CreateServiceOrder(g *gin.Context) {
	var serviceOrder entities.ServiceOrder
	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.serviceOrderUseCase.CreateServiceOrder(g.Request.Context(), serviceOrder)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceOrderNotFound) ||
			errors.Is(err, usecase.ErrVehicleNotFound) ||
			errors.Is(err, usecase.ErrCustomerNotFound) {
			g.JSON(404, gin.H{"error": err.Error()})
			return
		}

		g.JSON(500, gin.H{"error": "Failed to create service order", "details": err.Error()})
		log.Error().Msgf("Error creating service order: %v", err)
		return
	}

	g.JSON(201, result)
}

// UpdateServiceOrderDiagnosis godoc
// @Summary Update service order diagnosis
// @Description Update the diagnosis information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body entities.ServiceOrder true "Service Order Diagnosis Information"
// @Success 200 {object} entities.ServiceOrder
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/diagnosis [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderDiagnosis(g *gin.Context) {
	var serviceOrder entities.ServiceOrder

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid service order ID"})
		return
	}
	serviceOrder.ID = uint(id)

	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		log.Error().Msgf("Error binding JSON: %v", err)
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, DIAGNOSIS)
	if err != nil {

		if errors.Is(err, usecase.ErrServiceOrderNotFound) {
			g.JSON(404, gin.H{"error": "Service order not found"})
			return
		}

		if errors.Is(err, usecase.ErrInvalidTransitionStatusToDiagnosis) ||
			errors.Is(err, usecase.ErrInvalidStatus) ||
			errors.Is(err, usecase.ErrInvalidFlow) {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}

		g.JSON(500, gin.H{"error": "Failed to update service order", "details": err.Error()})
		return
	}

	g.JSON(200, result)
}

// UpdateServiceOrderEstimate godoc
// @Summary Update service order estimate
// @Description Update the estimate information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body entities.ServiceOrder true "Service Order Estimate Information"
// @Success 200 {object} entities.ServiceOrder
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/estimate [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderEstimate(g *gin.Context) {
	var serviceOrder entities.ServiceOrder

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid service order ID"})
		return
	}
	serviceOrder.ID = uint(id)

	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		log.Error().Msgf("Error binding JSON: %v", err)
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, ESTIMATE)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceOrderNotFound) {
			g.JSON(404, gin.H{"error": "Service order not found"})
			return
		}

		if errors.Is(err, usecase.ErrInvalidTransitionStatusToEstimate) ||
			errors.Is(err, usecase.ErrInvalidStatus) ||
			errors.Is(err, usecase.ErrInvalidFlow) {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}

		g.JSON(500, gin.H{"error": "Failed to update service order", "details": err.Error()})
		return
	}

	g.JSON(200, result)
}

// UpdateServiceOrderExecution godoc
// @Summary Update service order execution
// @Description Update the execution information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body entities.ServiceOrder true "Service Order Execution Information"
// @Success 200 {object} entities.ServiceOrder
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/execution [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderExecution(g *gin.Context) {
	var serviceOrder entities.ServiceOrder

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil || id <= 0 {
		g.JSON(400, gin.H{"error": "Invalid service order ID"})
		return
	}
	serviceOrder.ID = uint(id)

	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		log.Error().Msgf("Error binding JSON: %v", err)
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, EXECUTION)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceOrderNotFound) {
			g.JSON(404, gin.H{"error": "Service order not found"})
			return
		}

		if errors.Is(err, usecase.ErrInvalidTransitionStatusToExecution) ||
			errors.Is(err, usecase.ErrInvalidStatus) ||
			errors.Is(err, usecase.ErrInvalidFlow) {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}

		g.JSON(500, gin.H{"error": "Failed to update service order", "details": err.Error()})
		return
	}

	g.JSON(200, result)
}

// UpdateServiceOrderDelivery godoc
// @Summary Update service order delivery
// @Description Update the delivery information of a service order
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Param order body entities.ServiceOrder true "Service Order Delivery Information"
// @Success 200 {object} entities.ServiceOrder
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id}/delivery [patch]
func (h *ServiceOrderHandler) UpdateServiceOrderDelivery(g *gin.Context) {
	var serviceOrder entities.ServiceOrder

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil || id <= 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service order ID"})
		return
	}
	serviceOrder.ID = uint(id)

	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		log.Error().Msgf("Error binding JSON: %v", err)
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result, err := h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, DELIVERY)
	if err != nil {
		if errors.Is(err, usecase.ErrServiceOrderNotFound) {
			g.JSON(404, gin.H{"error": "Service order not found"})
			return
		}

		if errors.Is(err, usecase.ErrInvalidTransitionStatusToDelivery) ||
			errors.Is(err, usecase.ErrInvalidStatus) ||
			errors.Is(err, usecase.ErrInvalidFlow) {
			g.JSON(400, gin.H{"error": err.Error()})
			return
		}
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service order", "details": err.Error()})
		return
	}

	g.JSON(http.StatusOK, result)
}

// GetServiceOrder godoc
// @Summary Get service order by ID
// @Description Retrieve a service order by its ID
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Service Order ID"
// @Success 200 {object} entities.ServiceOrder
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /service-orders/{id} [get]
func (h *ServiceOrderHandler) GetServiceOrder(g *gin.Context) {
	var serviceOrder entities.ServiceOrder

	id, err := strconv.Atoi(g.Param("id"))
	if err != nil || id <= 0 {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service order ID"})
		return
	}
	serviceOrder.ID = uint(id)

	ServiceOrderResponse, err := h.serviceOrderUseCase.GetServiceOrder(g.Request.Context(), serviceOrder)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service order", "details": err.Error()})
		return
	}

	g.JSON(http.StatusOK, ServiceOrderResponse)
}

// ListServiceOrders godoc
// @Summary List all service orders
// @Description Get a list of all service orders
// @Tags Service Orders
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} entities.ServiceOrder
// @Failure 500 {object} map[string]string
// @Router /service-orders [get]
func (h *ServiceOrderHandler) ListServiceOrders(g *gin.Context) {
	serviceOrders, err := h.serviceOrderUseCase.ListServiceOrders(g.Request.Context())
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service orders", "details": err.Error()})
		return
	}

	g.JSON(http.StatusOK, serviceOrders)
}

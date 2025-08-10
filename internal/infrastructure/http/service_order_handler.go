package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
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

type ServiceOrderHandler struct {
	serviceOrderUseCase usecase.ServiceOrderUseCase
}

func NewServiceOrderHandler(useCase *usecase.ServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		serviceOrderUseCase: *useCase,
	}
}

// CreateServiceOrder POST /os
func (h *ServiceOrderHandler) CreateServiceOrder(g *gin.Context) {
	var serviceOrder entities.ServiceOrder
	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := h.serviceOrderUseCase.CreateServiceOrder(g.Request.Context(), serviceOrder)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to create service order"})
		return
	}

	g.JSON(201, gin.H{"message": "Service order created successfully"})
}

// UpdateServiceOrderDiagnosis PATCH /os/:id/diagnosis
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

	err = h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, DIAGNOSIS)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update service order"})
		return
	}

	g.JSON(200, gin.H{"message": "Service order updated successfully"})
}

// UpdateServiceOrderEstimate PATCH /os/:id/estimate
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

	err = h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, ESTIMATE)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update service order"})
		return
	}

	g.JSON(200, gin.H{"message": "Service order updated successfully"})
}

// UpdateServiceOrderExecution PATCH /os/:id/execution
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

	err = h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, EXECUTION)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update service order"})
		return
	}

	g.JSON(200, gin.H{"message": "Service order updated successfully"})
}

// UpdateServiceOrderDelivery PATCH /os/:id/delivery
func (h *ServiceOrderHandler) UpdateServiceOrderDelivery(g *gin.Context) {
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

	err = h.serviceOrderUseCase.UpdateServiceOrder(g.Request.Context(), serviceOrder, DELIVERY)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update service order"})
		return
	}

	g.JSON(200, gin.H{"message": "Service order updated successfully"})
}

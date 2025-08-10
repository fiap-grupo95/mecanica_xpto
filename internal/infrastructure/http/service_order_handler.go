package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type ServiceOrderHandler struct {
	UseCase usecase.ServiceOrderUseCase
}

func NewServiceOrderHandler(useCase *usecase.ServiceOrderUseCase) *ServiceOrderHandler {
	return &ServiceOrderHandler{
		UseCase: *useCase,
	}
}

func (h *ServiceOrderHandler) CreateServiceOrder(g *gin.Context) {
	var serviceOrder entities.ServiceOrder
	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := h.UseCase.CreateServiceOrder(serviceOrder)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to create service order"})
		return
	}

	g.JSON(201, gin.H{"message": "Service order created successfully"})
}

func (h *ServiceOrderHandler) UpdateServiceOrder(g *gin.Context) {
	var serviceOrder entities.ServiceOrder
	if err := g.ShouldBindJSON(&serviceOrder); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := h.UseCase.UpdateServiceOrder(serviceOrder)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update service order"})
		return
	}

	g.JSON(200, gin.H{"message": "Service order updated successfully"})
}

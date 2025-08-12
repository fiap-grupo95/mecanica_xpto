package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdditionalRepairHandler struct {
	additionalRepairUseCase usecase.IAdditionalRepairUseCase
}

func NewAdditionalRepairHandler(useCase usecase.IAdditionalRepairUseCase) *AdditionalRepairHandler {
	return &AdditionalRepairHandler{
		additionalRepairUseCase: useCase,
	}
}

func (h *AdditionalRepairHandler) GetAdditionalRepair(g *gin.Context) {
	idStr := g.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		appErr := pkg.NewDomainErrorSimple("INVALID_ID", "Invalid additional repair ID", http.StatusBadRequest)
		g.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
	}

	foundAdr, err := h.additionalRepairUseCase.GetAdditionalRepair(g.Request.Context(), uint(id))
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to get additional repair"})
		return
	}

	g.JSON(http.StatusOK, foundAdr)
}

// CreateSOAdditionalRepair POST /os
func (h *AdditionalRepairHandler) CreateSOAdditionalRepair(g *gin.Context) {
	var adr entities.AdditionalRepair
	if err := g.ShouldBindJSON(&adr); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err := h.additionalRepairUseCase.CreateAdditionalRepair(g.Request.Context(), adr)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to create service order"})
		return
	}

	g.JSON(201, gin.H{"message": "Service order created successfully"})
}

func (h *AdditionalRepairHandler) AddPartSupplyAndService(g *gin.Context) {
	idStr := g.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		appErr := pkg.NewDomainErrorSimple("INVALID_ID", "Invalid additional repair ID", http.StatusBadRequest)
		g.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
	}
	var adr entities.AdditionalRepair
	if err := g.ShouldBindJSON(&adr); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err = h.additionalRepairUseCase.AddPartSupplyAndService(g.Request.Context(), uint(id), adr)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update additional repair"})
		return
	}

	g.JSON(201, gin.H{"message": "Additional repair updated successfully"})
}

func (h *AdditionalRepairHandler) RemovePartSupplyAndService(g *gin.Context) {
	idStr := g.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		appErr := pkg.NewDomainErrorSimple("INVALID_ID", "Invalid customer ID", http.StatusBadRequest)
		g.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
	}
	var adr entities.AdditionalRepair
	if err := g.ShouldBindJSON(&adr); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err = h.additionalRepairUseCase.RemovePartSupplyAndService(g.Request.Context(), uint(id), adr)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update additional repair"})
		return
	}

	g.JSON(201, gin.H{"message": "Additional repair updated successfully"})
}

func (h *AdditionalRepairHandler) CustomerApproval(g *gin.Context) {
	idStr := g.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		appErr := pkg.NewDomainErrorSimple("INVALID_ID", "Invalid customer ID", http.StatusBadRequest)
		g.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
	}
	var adr entities.AdditionalRepairStatusDTO
	if err := g.ShouldBindJSON(&adr); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	err = h.additionalRepairUseCase.CustomerApprovalStatus(g.Request.Context(), uint(id), adr)
	if err != nil {
		g.JSON(500, gin.H{"error": "Failed to update additional repair"})
		return
	}

	g.JSON(201, gin.H{"message": "Additional repair updated successfully"})
}

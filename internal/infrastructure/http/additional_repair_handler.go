package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdditionalRepairHandler handles HTTP requests for additional repairs
// @title Additional Repair API
// @version 1.0
// @description API for managing additional repairs in the workshop management system
type AdditionalRepairHandler struct {
	additionalRepairUseCase usecase.IAdditionalRepairUseCase
}

func NewAdditionalRepairHandler(useCase usecase.IAdditionalRepairUseCase) *AdditionalRepairHandler {
	return &AdditionalRepairHandler{
		additionalRepairUseCase: useCase,
	}
}

// GetAdditionalRepair godoc
// @Summary Get additional repair by ID
// @Description Retrieve an additional repair by its ID
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Success 200 {object} entities.AdditionalRepair
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id} [get]
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

// CreateAdditionalRepair godoc
// @Summary Create a new additional repair
// @Description Create a new additional repair record
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param repair body entities.AdditionalRepair true "Additional Repair Information"
// @Success 201 {object} entities.AdditionalRepair
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs [post]
func (h *AdditionalRepairHandler) CreateAdditionalRepair(g *gin.Context) {
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

// AddPartSupplyAndService godoc
// @Summary Add parts supply and service to additional repair
// @Description Add parts supply and service to an existing additional repair
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Param repair body entities.AdditionalRepair true "Parts Supply and Service Information"
// @Success 201 {object} map[string]string
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id}/add [post]
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

// RemovePartSupplyAndService godoc
// @Summary Remove parts supply and service from additional repair
// @Description Remove parts supply and service from an existing additional repair
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Param repair body entities.AdditionalRepair true "Parts Supply and Service Information"
// @Success 201 {object} map[string]string
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id}/remove [delete]
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

// CustomerApproval godoc
// @Summary Update customer approval status for additional repair
// @Description Update the customer approval status of an additional repair
// @Tags Additional Repairs
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Additional Repair ID"
// @Param status body entities.AdditionalRepairStatusDTO true "Approval Status Information"
// @Success 201 {object} map[string]string
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /additional-repairs/{id}/customer_approval [post]
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

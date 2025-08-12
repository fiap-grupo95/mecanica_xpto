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
	errInvalidPartsSupplyID    = pkg.NewDomainErrorSimple("INVALID_PARTS_SUPPLY_ID", "Invalid parts supply ID", http.StatusBadRequest)
	errInvalidPartsSupplyInput = pkg.NewDomainErrorSimple("INVALID_INPUT", "Invalid input data", http.StatusBadRequest)
)

// PartsSupplyHandler handles HTTP requests for parts supply operations
// @title Parts Supply API
// @version 1.0
// @description API for managing parts supply in the workshop management system
type PartsSupplyHandler struct {
	usecase usecase.IPartsSupplyUseCase
}

func NewPartsSupplyHandler(usecase usecase.IPartsSupplyUseCase) *PartsSupplyHandler {
	return &PartsSupplyHandler{usecase: usecase}
}

func mapPartsSupplyError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrPartsSupplyNotFound):
		return pkg.NewDomainErrorSimple("PARTS_SUPPLY_NOT_FOUND", "parts supply not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid parts supply ID", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrPartsSupplyAlreadyExists):
		return pkg.NewDomainErrorSimple("PARTS_SUPPLY_EXISTS", "parts supply already exists", http.StatusConflict)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetPartsSupplyByID godoc
// @Summary Get parts supply by ID
// @Description Retrieve a parts supply by its ID
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Parts Supply ID"
// @Success 200 {object} entities.PartsSupply
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies/{id} [get]
func (h *PartsSupplyHandler) GetPartsSupplyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(errInvalidPartsSupplyID.HTTPStatus, errInvalidPartsSupplyID.ToHTTPError())
		return
	}

	foundPartsSupply, err := h.usecase.GetPartsSupplyByID(c.Request.Context(), uint(id))
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, foundPartsSupply)
}

// CreatePartsSupply godoc
// @Summary Create a new parts supply
// @Description Create a new parts supply record
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param supply body entities.PartsSupply true "Parts Supply Information"
// @Success 201 {object} entities.PartsSupply
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies [post]
func (h *PartsSupplyHandler) CreatePartsSupply(c *gin.Context) {
	var partsSupply entities.PartsSupply
	if err := c.ShouldBindJSON(&partsSupply); err != nil {
		c.JSON(errInvalidPartsSupplyInput.HTTPStatus, errInvalidPartsSupplyInput.ToHTTPError())
		return
	}

	createdPartsSupply, err := h.usecase.CreatePartsSupply(c.Request.Context(), &partsSupply)
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, createdPartsSupply)
}

// UpdatePartsSupply godoc
// @Summary Update a parts supply
// @Description Update an existing parts supply record
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Parts Supply ID"
// @Param supply body entities.PartsSupply true "Parts Supply Information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies/{id} [put]
func (h *PartsSupplyHandler) UpdatePartsSupply(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var partsSupply entities.PartsSupply
	if err := c.ShouldBindJSON(&partsSupply); err != nil {
		c.JSON(errInvalidPartsSupplyInput.HTTPStatus, errInvalidPartsSupplyInput.ToHTTPError())
		return
	}
	partsSupply.ID = uint(id)

	if err := h.usecase.UpdatePartsSupply(c.Request.Context(), &partsSupply); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "parts supply updated successfully"})
}

// DeletePartsSupply godoc
// @Summary Delete a parts supply
// @Description Delete an existing parts supply record
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Parts Supply ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies/{id} [delete]
func (h *PartsSupplyHandler) DeletePartsSupply(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(errInvalidPartsSupplyID.HTTPStatus, errInvalidPartsSupplyID.ToHTTPError())
		return
	}

	if err := h.usecase.DeletePartsSupply(c.Request.Context(), uint(id)); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "parts supply deleted successfully"})
}

// ListPartsSupplies godoc
// @Summary List all parts supplies
// @Description Get a list of all parts supplies
// @Tags Parts Supply
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} entities.PartsSupply
// @Failure 500 {object} pkg.AppError
// @Router /parts-supplies [get]
func (h *PartsSupplyHandler) ListPartsSupplies(c *gin.Context) {
	partsSupplies, err := h.usecase.ListPartsSupplies(c.Request.Context())
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, partsSupplies)
}

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
	errInvalidPartsSupplyID = pkg.NewDomainErrorSimple("INVALID_PARTS_SUPPLY_ID", "Invalid parts supply ID", http.StatusBadRequest)
)

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

func (h *PartsSupplyHandler) CreatePartsSupply(c *gin.Context) {
	var partsSupply entities.PartsSupply
	if err := c.ShouldBindJSON(&partsSupply); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
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

func (h *PartsSupplyHandler) UpdatePartsSupply(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var partsSupply entities.PartsSupply
	if err := c.ShouldBindJSON(&partsSupply); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
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

func (h *PartsSupplyHandler) DeletePartsSupply(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(errInvalidPartsSupplyID.HTTPStatus, errInvalidPartsSupplyID.ToHTTPError)
		return
	}

	if err := h.usecase.DeletePartsSupply(c.Request.Context(), uint(id)); err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "parts supply deleted successfully"})
}

func (h *PartsSupplyHandler) ListPartsSupplies(c *gin.Context) {
	partsSupplies, err := h.usecase.ListPartsSupplies(c.Request.Context())
	if err != nil {
		appErr := mapPartsSupplyError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, partsSupplies)
}

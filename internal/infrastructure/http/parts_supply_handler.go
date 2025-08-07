package http

import (
	"mecanica_xpto/internal/domain/model/entities"
	use_case "mecanica_xpto/internal/domain/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PartsSupplyHandler struct {
	usecase use_case.IPartsSupplyUseCase
}

func NewPartsSupplyHandler(usecase use_case.IPartsSupplyUseCase) *PartsSupplyHandler {
	return &PartsSupplyHandler{usecase: usecase}
}

func (h *PartsSupplyHandler) GetPartsSupplyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	foundPartsSupply, err := h.usecase.GetPartsSupplyByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == use_case.ErrPartsSupplyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "parts supply not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, foundPartsSupply)
}

func (h *PartsSupplyHandler) CreatePartsSupply(c *gin.Context) {
	var partsSupply entities.PartsSupply
	if err := c.ShouldBindJSON(&partsSupply); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	createdPartsSupply, err := h.usecase.CreatePartsSupply(c.Request.Context(), &partsSupply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create parts supply"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	partsSupply.ID = uint(id)

	if err := h.usecase.UpdatePartsSupply(c.Request.Context(), &partsSupply); err != nil {
		if err == use_case.ErrPartsSupplyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "parts supply not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update parts supply"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "parts supply updated successfully"})
}

func (h *PartsSupplyHandler) DeletePartsSupply(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.usecase.DeletePartsSupply(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete parts supply"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "parts supply deleted successfully"})
}

func (h *PartsSupplyHandler) ListPartsSupplies(c *gin.Context) {
	partsSupplies, err := h.usecase.ListPartsSupplies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve parts supplies"})
		return
	}

	c.JSON(http.StatusOK, partsSupplies)
}

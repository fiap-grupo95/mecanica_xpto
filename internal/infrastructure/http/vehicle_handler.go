package http

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/pkg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	// Domain errors
	errInvalidVehicleID = pkg.NewDomainErrorSimple("INVALID_VEHICLE_ID", "Invalid vehicle ID", http.StatusBadRequest)
	errInvalidInput     = pkg.NewDomainErrorSimple("INVALID_INPUT", "Invalid input data", http.StatusBadRequest)
	errEmptyPlate       = pkg.NewDomainErrorSimple("EMPTY_PLATE", "Plate cannot be empty", http.StatusBadRequest)
)

// VehicleHandler handles HTTP requests for vehicle operations
// @title Vehicle API
type VehicleHandler struct {
	service usecase.VehicleServiceInterface
}

func NewVehicleHandler(service usecase.VehicleServiceInterface) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

func mapVehicleError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrVehicleNotFound):
		return pkg.NewDomainErrorSimple("VEHICLE_NOT_FOUND", "Vehicle not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrInvalidPlateFormat):
		return pkg.NewDomainErrorSimple("INVALID_PLATE_FORMAT", "Invalid plate format", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrVehicleAlreadyExists):
		return pkg.NewDomainErrorSimple("VEHICLE_EXISTS", "Vehicle already exists", http.StatusConflict)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid vehicle ID", http.StatusBadRequest)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetVehicles godoc
// @Summary Get all vehicles
// @Description Retrieves a list of all vehicles
// @Tags vehicles
// @Accept json
// @Produce json
// @Success 200 {array} entities.Vehicle
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles [get]
func (v VehicleHandler) GetVehicles(c *gin.Context) {
	vehicles, err := v.service.GetAllVehicles()
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

// GetVehicleByCustomerID godoc
// @Summary Get vehicles by customer ID
// @Description Retrieves all vehicles belonging to a specific customer
// @Tags vehicles
// @Accept json
// @Produce json
// @Param customerID path int true "Customer ID"
// @Success 200 {array} entities.Vehicle
// @Failure 400 {object} pkg.ErrorResponse "Invalid customer ID format"
// @Failure 404 {object} pkg.ErrorResponse "Customer not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/customer/{customerID} [get]
func (v VehicleHandler) GetVehicleByCustomerID(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("customerID"), 10, 32)
	if err != nil {
		c.JSON(errInvalidVehicleID.HTTPStatus, errInvalidVehicleID.ToHTTPError())
		return
	}

	vehicles, err := v.service.GetVehiclesByCustomerID(uint(customerID))
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, vehicles)
}

// GetVehicleByID godoc
// @Summary Get vehicle by ID
// @Description Retrieves a vehicle by its ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 200 {object} entities.Vehicle
// @Failure 400 {object} pkg.ErrorResponse "Invalid vehicle ID format"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/{id} [get]
func (v VehicleHandler) GetVehicleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(errInvalidVehicleID.HTTPStatus, errInvalidVehicleID.ToHTTPError())
		return
	}

	vehicle, err := v.service.GetVehicleByID(uint(id))
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// GetVehicleByPlate godoc
// @Summary Get vehicle by plate
// @Description Retrieves a vehicle by its license plate
// @Tags vehicles
// @Accept json
// @Produce json
// @Param plate path string true "Vehicle license plate"
// @Success 200 {object} entities.Vehicle
// @Failure 400 {object} pkg.ErrorResponse "Invalid plate format or empty plate"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/plate/{plate} [get]
func (v VehicleHandler) GetVehicleByPlate(c *gin.Context) {
	plate := c.Param("plate")
	if plate == "" {
		c.JSON(errEmptyPlate.HTTPStatus, errEmptyPlate.ToHTTPError())
		return
	}

	vehicle, err := v.service.GetVehicleByPlate(plate)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, vehicle)
}

// CreateVehicle godoc
// @Summary Create a new vehicle
// @Description Creates a new vehicle record
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vehicle body entities.Vehicle true "Vehicle information"
// @Success 201 {object} pkg.ErrorResponse{message=string} "Vehicle created successfully"
// @Failure 400 {object} pkg.ErrorResponse "Invalid input data or plate format"
// @Failure 409 {object} pkg.ErrorResponse "Vehicle already exists"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles [post]
func (v VehicleHandler) CreateVehicle(c *gin.Context) {
	var vehicle entities.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
		return
	}

	result, err := v.service.CreateVehicle(vehicle)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": result})
}

// UpdateVehicle godoc
// @Summary Update a vehicle partially
// @Description Updates specific fields of an existing vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} pkg.ErrorResponse{message=string} "Vehicle updated successfully"
// @Failure 400 {object} pkg.ErrorResponse "Invalid input data, ID format or plate format"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/{id} [patch]
func (v VehicleHandler) UpdateVehicle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(errInvalidVehicleID.HTTPStatus, errInvalidVehicleID.ToHTTPError())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(errInvalidInput.HTTPStatus, errInvalidInput.ToHTTPError())
		return
	}

	result, err := v.service.UpdateVehiclePartial(uint(id), updates)
	if err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": result})
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Deletes a vehicle by its ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Success 204 "No Content"
// @Failure 400 {object} pkg.ErrorResponse "Invalid vehicle ID format"
// @Failure 404 {object} pkg.ErrorResponse "Vehicle not found"
// @Failure 500 {object} pkg.ErrorResponse "Internal server error"
// @Router /vehicles/{id} [delete]
func (v VehicleHandler) DeleteVehicle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(errInvalidVehicleID.HTTPStatus, errInvalidVehicleID.ToHTTPError())
		return
	}

	if err := v.service.DeleteVehicle(uint(id)); err != nil {
		appErr := mapVehicleError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.Status(http.StatusNoContent)
}

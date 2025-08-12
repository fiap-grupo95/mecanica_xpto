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
	errInvalidPaymentID    = pkg.NewDomainErrorSimple("INVALID_PAYMENT_ID", "Invalid payment ID", http.StatusBadRequest)
	errInvalidPaymentInput = pkg.NewDomainErrorSimple("INVALID_PAYMENT_INPUT", "Invalid payment input", http.StatusBadRequest)
)

// PaymentHandler handles HTTP requests for payment operations
// @title Payment API
// @version 1.0
// @description API for managing payments in the workshop management system
type PaymentHandler struct {
	usecase usecase.IPaymentUseCase
}

func NewPaymentHandler(usecase usecase.IPaymentUseCase) *PaymentHandler {
	return &PaymentHandler{usecase: usecase}
}

func mapPaymentError(err error) *pkg.AppError {
	switch {
	case errors.Is(err, usecase.ErrorPaymentNotFound):
		return pkg.NewDomainErrorSimple("PAYMENT_NOT_FOUND", "Payment not found", http.StatusNotFound)
	case errors.Is(err, usecase.ErrPaymentAmountDoesNotMatch):
		return pkg.NewDomainErrorSimple("PAYMENT_AMOUNT_DOES_NOT_MATCH", "Payment amount does not match service order estimate", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrInvalidID):
		return pkg.NewDomainErrorSimple("INVALID_ID", "Invalid payment ID", http.StatusBadRequest)
	case errors.Is(err, usecase.ErrPaymentAlreadyExists):
		return pkg.NewDomainErrorSimple("PAYMENT_ALREADY_EXISTS", "Payment already exists", http.StatusConflict)
	default:
		return pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", err, http.StatusInternalServerError)
	}
}

// GetPaymentByID godoc
// @Summary Get payment by ID
// @Description Retrieve a payment by its ID
// @Tags Payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} entities.Payment
// @Failure 400 {object} pkg.AppError
// @Failure 404 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(errInvalidPaymentID.HTTPStatus, errInvalidPaymentID.ToHTTPError())
		return
	}

	payment, err := h.usecase.GetPaymentByID(c.Request.Context(), uint(id))
	if err != nil {
		appErr := mapPaymentError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, payment)
}

// ListPayments godoc
// @Summary List all payments
// @Description Get a list of all payments
// @Tags Payments
// @Security Bearer
// @Accept json
// @Produce json
// @Success 200 {array} entities.Payment
// @Failure 500 {object} pkg.AppError
// @Router /payments [get]
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	payments, err := h.usecase.ListPayments(c.Request.Context())
	if err != nil {
		appErr := mapPaymentError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusOK, payments)
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a new payment record
// @Tags Payments
// @Security Bearer
// @Accept json
// @Produce json
// @Param payment body entities.Payment true "Payment Information"
// @Success 201 {object} entities.Payment
// @Failure 400 {object} pkg.AppError
// @Failure 500 {object} pkg.AppError
// @Router /payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var input entities.Payment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(errInvalidPaymentInput.HTTPStatus, errInvalidPaymentInput.ToHTTPError())
		return
	}
	payment, err := h.usecase.CreatePayment(c.Request.Context(), &input)
	if err != nil {
		appErr := mapPaymentError(err)
		c.JSON(appErr.HTTPStatus, appErr.ToHTTPError())
		return
	}

	c.JSON(http.StatusCreated, payment)
}

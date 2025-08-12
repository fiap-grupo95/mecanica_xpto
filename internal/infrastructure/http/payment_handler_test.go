package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"

	mocks "mecanica_xpto/internal/domain/mocks"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
)

func setupPaymentHandlerTest(t *testing.T) (*mocks.MockIPaymentUseCase, *PaymentHandler, *gin.Engine) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockIPaymentUseCase(ctrl)
	h := NewPaymentHandler(mockUC)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return mockUC, h, r
}

func TestPaymentHandler_GetPaymentByID(t *testing.T) {
	mockUC, h, r := setupPaymentHandlerTest(t)
	r.GET("/v1/payments/:id", h.GetPaymentByID)

	payment := &entities.Payment{ID: 1, ServiceOrderID: 1, Amount: 100.0}

	t.Run("success", func(t *testing.T) {
		mockUC.EXPECT().GetPaymentByID(gomock.Any(), uint(1)).Return(payment, nil)
		req, _ := http.NewRequest(http.MethodGet, "/v1/payments/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/payments/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockUC.EXPECT().GetPaymentByID(gomock.Any(), uint(2)).Return(&entities.Payment{}, usecase.ErrorPaymentNotFound)
		req, _ := http.NewRequest(http.MethodGet, "/v1/payments/2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected 404, got %d", w.Code)
		}
	})
}

func TestPaymentHandler_ListPayments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC, h, r := setupPaymentHandlerTest(t)
	r.GET("/v1/payments", h.ListPayments)

	payments := []entities.Payment{{ID: 1}, {ID: 2}}

	t.Run("success", func(t *testing.T) {
		mockUC.EXPECT().ListPayments(gomock.Any()).Return(payments, nil)
		req, _ := http.NewRequest(http.MethodGet, "/v1/payments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", w.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mockUC.EXPECT().ListPayments(gomock.Any()).Return(nil, errors.New("fail"))
		req, _ := http.NewRequest(http.MethodGet, "/v1/payments", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

func TestPaymentHandler_CreatePayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC, h, r := setupPaymentHandlerTest(t)
	r.POST("/v1/payments", h.CreatePayment)

	payment := &entities.Payment{ID: 1, ServiceOrderID: 1, Amount: 100.0}

	t.Run("success", func(t *testing.T) {
		mockUC.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Return(payment, nil)
		body, _ := json.Marshal(payment)
		req, _ := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d", w.Code)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("amount does not match", func(t *testing.T) {
		mockUC.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrPaymentAmountDoesNotMatch)
		body, _ := json.Marshal(payment)
		req, _ := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", w.Code)
		}
	})

	t.Run("already exists", func(t *testing.T) {
		mockUC.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Return(nil, usecase.ErrPaymentAlreadyExists)
		body, _ := json.Marshal(payment)
		req, _ := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", w.Code)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mockUC.EXPECT().CreatePayment(gomock.Any(), gomock.Any()).Return(nil, errors.New("fail"))
		body, _ := json.Marshal(payment)
		req, _ := http.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", w.Code)
		}
	})
}

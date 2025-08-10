package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/infrastructure/http"
	"mecanica_xpto/internal/infrastructure/http/mocks"
	h "net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *http.CustomerHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/customers/:document", handler.GetCustomer)
	r.GET("/customers/id/:id", handler.GetFullCustomer)
	r.POST("/customers", handler.CreateCustomer)
	r.PUT("/customers/:id", handler.UpdateCustomer)
	r.DELETE("/customers/:id", handler.DeleteCustomer)
	r.GET("/customers", handler.ListCustomer)
	return r
}

func TestGetCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := http.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	customer := &entities.Customer{FullName: "Test"}
	mockUC.EXPECT().GetByDocument("123").Return(customer, nil)

	req, _ := h.NewRequest("GET", "/customers/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, h.StatusOK, w.Code)
}

func TestGetCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := http.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().GetByDocument("123").Return(nil, errors.New("not found"))

	req, _ := h.NewRequest("GET", "/customers/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, h.StatusInternalServerError, w.Code)
}

func TestCreateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := http.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	customer := entities.Customer{FullName: "Test"}
	body, _ := json.Marshal(customer)
	mockUC.EXPECT().CreateCustomer(gomock.Any()).Return(nil)

	req, _ := h.NewRequest("POST", "/customers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, h.StatusCreated, w.Code)
}

func TestUpdateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := http.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	customer := entities.Customer{FullName: "Test"}
	body, _ := json.Marshal(customer)
	mockUC.EXPECT().UpdateCustomer(uint(1), gomock.Any()).Return(nil)

	req, _ := h.NewRequest("PUT", "/customers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, h.StatusOK, w.Code)
}

func TestDeleteCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := http.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().DeleteCustomer(uint(1)).Return(nil)

	req, _ := h.NewRequest("DELETE", "/customers/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, h.StatusNoContent, w.Code)
}

func TestListCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockICustomerUseCase(ctrl)
	handler := http.NewCustomerHandler(mockUC)
	router := setupRouter(handler)

	mockUC.EXPECT().ListCustomer().Return([]entities.Customer{}, nil)

	req, _ := h.NewRequest("GET", "/customers", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, h.StatusOK, w.Code)
}

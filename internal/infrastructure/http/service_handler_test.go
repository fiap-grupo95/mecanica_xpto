package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mecanica_xpto/internal/domain/mocks"
	"mecanica_xpto/internal/domain/model/entities"
	use_case "mecanica_xpto/internal/domain/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func setupServiceHandlerTest(t *testing.T) (*mocks.MockIServiceUseCase, *ServiceHandler, *gin.Engine) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockIServiceUseCase(ctrl)
	h := NewServiceHandler(mockUC)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return mockUC, h, r
}

func TestGetServiceByID(t *testing.T) {
	mockUC, h, r := setupServiceHandlerTest(t)
	r.GET("/service/:id", h.GetServiceByID)

	mockUC.EXPECT().GetServiceByID(gomock.Any(), uint(1)).Return(entities.Service{ID: 1, Name: "Troca de óleo"}, nil)
	req, _ := http.NewRequest("GET", "/service/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().GetServiceByID(gomock.Any(), uint(2)).Return(entities.Service{}, use_case.ErrServiceNotFound)
	req, _ = http.NewRequest("GET", "/service/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	mockUC.EXPECT().GetServiceByID(gomock.Any(), uint(3)).Return(entities.Service{}, errors.New("fail"))
	req, _ = http.NewRequest("GET", "/service/3", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = http.NewRequest("GET", "/service/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateService(t *testing.T) {
	mockUC, h, r := setupServiceHandlerTest(t)
	r.POST("/service", h.CreateService)
	jsonBody := `{"name":"Troca de óleo"}`

	mockUC.EXPECT().CreateService(gomock.Any(), gomock.Any()).Return(entities.Service{ID: 1, Name: "Troca de óleo"}, nil)
	req, _ := http.NewRequest("POST", "/service", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	mockUC.EXPECT().CreateService(gomock.Any(), gomock.Any()).Return(entities.Service{}, errors.New("fail"))
	req, _ = http.NewRequest("POST", "/service", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = http.NewRequest("POST", "/service", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateService(t *testing.T) {
	mockUC, h, r := setupServiceHandlerTest(t)
	r.PUT("/service/:id", h.UpdateService)
	jsonBody := `{"name":"Troca de óleo"}`

	mockUC.EXPECT().UpdateService(gomock.Any(), gomock.Any()).Return(nil)
	req, _ := http.NewRequest("PUT", "/service/1", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().UpdateService(gomock.Any(), gomock.Any()).Return(use_case.ErrServiceNotFound)
	req, _ = http.NewRequest("PUT", "/service/2", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	mockUC.EXPECT().UpdateService(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
	req, _ = http.NewRequest("PUT", "/service/3", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = http.NewRequest("PUT", "/service/abc", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	req, _ = http.NewRequest("PUT", "/service/1", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeleteService(t *testing.T) {
	mockUC, h, r := setupServiceHandlerTest(t)
	r.DELETE("/service/:id", h.DeleteService)

	mockUC.EXPECT().DeleteService(gomock.Any(), uint(1)).Return(nil)
	req, _ := http.NewRequest("DELETE", "/service/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}

	mockUC.EXPECT().DeleteService(gomock.Any(), uint(2)).Return(use_case.ErrServiceNotFound)
	req, _ = http.NewRequest("DELETE", "/service/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	mockUC.EXPECT().DeleteService(gomock.Any(), uint(3)).Return(errors.New("fail"))
	req, _ = http.NewRequest("DELETE", "/service/3", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = http.NewRequest("DELETE", "/service/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListServices(t *testing.T) {
	mockUC, h, r := setupServiceHandlerTest(t)
	r.GET("/service", h.ListServices)
	services := []entities.Service{{ID: 1, Name: "Troca de óleo"}, {ID: 2, Name: "Alinhamento"}}

	mockUC.EXPECT().ListServices(gomock.Any()).Return(services, nil)
	req, _ := http.NewRequest("GET", "/service", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().ListServices(gomock.Any()).Return(nil, errors.New("fail"))
	req, _ = http.NewRequest("GET", "/service", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

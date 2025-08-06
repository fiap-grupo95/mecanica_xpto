package http_test

import (
	"bytes"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"mecanica_xpto/internal/domain/mocks"
	"mecanica_xpto/internal/domain/model/entities"
	httpapi "mecanica_xpto/internal/infrastructure/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func setupPartsSupplyHandlerTest(t *testing.T) (*mocks.MockIPartsSupplyUseCase, *httpapi.PartsSupplyHandler, *gin.Engine) {
	ctrl := gomock.NewController(t)
	mockUC := mocks.NewMockIPartsSupplyUseCase(ctrl)
	h := httpapi.NewPartsSupplyHandler(mockUC)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return mockUC, h, r
}

func TestGetPartsSupplyByID(t *testing.T) {
	mockUC, h, r := setupPartsSupplyHandlerTest(t)
	r.GET("/parts/:id", h.GetPartsSupplyByID)

	mockUC.EXPECT().GetPartsSupplyByID(gomock.Any(), uint(1)).Return(entities.PartsSupply{ID: 1, Name: "Filtro"}, nil)
	req, _ := stdhttp.NewRequest("GET", "/parts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().GetPartsSupplyByID(gomock.Any(), uint(2)).Return(entities.PartsSupply{}, errors.New("fail"))
	req, _ = stdhttp.NewRequest("GET", "/parts/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = stdhttp.NewRequest("GET", "/parts/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreatePartsSupply(t *testing.T) {
	mockUC, h, r := setupPartsSupplyHandlerTest(t)
	r.POST("/parts", h.CreatePartsSupply)
	jsonBody := `{"name":"Filtro"}`

	mockUC.EXPECT().CreatePartsSupply(gomock.Any(), gomock.Any()).Return(entities.PartsSupply{ID: 1, Name: "Filtro"}, nil)
	req, _ := stdhttp.NewRequest("POST", "/parts", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	mockUC.EXPECT().CreatePartsSupply(gomock.Any(), gomock.Any()).Return(entities.PartsSupply{}, errors.New("fail"))
	req, _ = stdhttp.NewRequest("POST", "/parts", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = stdhttp.NewRequest("POST", "/parts", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdatePartsSupply(t *testing.T) {
	mockUC, h, r := setupPartsSupplyHandlerTest(t)
	r.PUT("/parts/:id", h.UpdatePartsSupply)
	jsonBody := `{"name":"Filtro"}`

	mockUC.EXPECT().UpdatePartsSupply(gomock.Any(), gomock.Any()).Return(nil)
	req, _ := stdhttp.NewRequest("PUT", "/parts/1", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().UpdatePartsSupply(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
	req, _ = stdhttp.NewRequest("PUT", "/parts/2", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = stdhttp.NewRequest("PUT", "/parts/abc", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	req, _ = stdhttp.NewRequest("PUT", "/parts/1", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestDeletePartsSupply(t *testing.T) {
	mockUC, h, r := setupPartsSupplyHandlerTest(t)
	r.DELETE("/parts/:id", h.DeletePartsSupply)

	mockUC.EXPECT().DeletePartsSupply(gomock.Any(), uint(1)).Return(nil)
	req, _ := stdhttp.NewRequest("DELETE", "/parts/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().DeletePartsSupply(gomock.Any(), uint(2)).Return(errors.New("fail"))
	req, _ = stdhttp.NewRequest("DELETE", "/parts/2", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	req, _ = stdhttp.NewRequest("DELETE", "/parts/abc", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListPartsSupplies(t *testing.T) {
	mockUC, h, r := setupPartsSupplyHandlerTest(t)
	r.GET("/parts", h.ListPartsSupplies)
	parts := []entities.PartsSupply{{ID: 1, Name: "Filtro"}, {ID: 2, Name: "Pastilha"}}

	mockUC.EXPECT().ListPartsSupplies(gomock.Any()).Return(parts, nil)
	req, _ := stdhttp.NewRequest("GET", "/parts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	mockUC.EXPECT().ListPartsSupplies(gomock.Any()).Return(nil, errors.New("fail"))
	req, _ = stdhttp.NewRequest("GET", "/parts", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != stdhttp.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

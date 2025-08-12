package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mecanica_xpto/internal/domain/model/entities"
	handler "mecanica_xpto/internal/infrastructure/http"
	"mecanica_xpto/internal/infrastructure/http/mocks"
)

func setupADRRouter(h *handler.AdditionalRepairHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/additional-repair", h.CreateSOAdditionalRepair)
	r.GET("/additional-repair/:id", h.GetAdditionalRepair)
	r.POST("/additional-repair/:id/part", h.AddPartSupplyAndService)
	r.DELETE("/additional-repair/:id/part", h.RemovePartSupplyAndService)
	r.POST("/additional-repair/:id/approval", h.CustomerApproval)
	return r
}

func TestCreateSOAdditionalRepair_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	adr := entities.AdditionalRepair{ServiceOrderID: 1, Description: "desc"}
	mockUC.EXPECT().CreateAdditionalRepair(gomock.Any(), adr).Return(nil)

	body, _ := json.Marshal(adr)
	req, _ := http.NewRequest("POST", "/additional-repair", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateSOAdditionalRepair_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	adr := entities.AdditionalRepair{ServiceOrderID: 1, Description: "desc"}
	mockUC.EXPECT().CreateAdditionalRepair(gomock.Any(), adr).Return(errors.New("fail"))

	body, _ := json.Marshal(adr)
	req, _ := http.NewRequest("POST", "/additional-repair", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// Go
func TestGetAdditionalRepair_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	expected := entities.AdditionalRepair{ID: 1, Description: "desc"}
	mockUC.EXPECT().GetAdditionalRepair(gomock.Any(), uint(1)).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/additional-repair/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entities.AdditionalRepair
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
}

func TestGetAdditionalRepair_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	// Return empty struct and error
	mockUC.EXPECT().GetAdditionalRepair(gomock.Any(), uint(999)).Return(entities.AdditionalRepair{}, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/additional-repair/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestAddPartSupplyAndService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	adr := entities.AdditionalRepair{ID: 1}
	mockUC.EXPECT().AddPartSupplyAndService(gomock.Any(), uint(1), adr).Return(nil)

	body, _ := json.Marshal(adr)
	req, _ := http.NewRequest("POST", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddPartSupplyAndService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	adr := entities.AdditionalRepair{ID: 1}
	mockUC.EXPECT().AddPartSupplyAndService(gomock.Any(), uint(1), adr).Return(errors.New("fail"))

	body, _ := json.Marshal(adr)
	req, _ := http.NewRequest("POST", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRemovePartSupplyAndService_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	adr := entities.AdditionalRepair{ID: 1}
	mockUC.EXPECT().RemovePartSupplyAndService(gomock.Any(), uint(1), adr).Return(nil)

	body, _ := json.Marshal(adr)
	req, _ := http.NewRequest("DELETE", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRemovePartSupplyAndService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	adr := entities.AdditionalRepair{ID: 1}
	mockUC.EXPECT().RemovePartSupplyAndService(gomock.Any(), uint(1), adr).Return(errors.New("fail"))

	body, _ := json.Marshal(adr)
	req, _ := http.NewRequest("DELETE", "/additional-repair/1/part", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCustomerApproval_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	dto := entities.AdditionalRepairStatusDTO{ApprovalStatus: "APPROVED"}
	mockUC.EXPECT().CustomerApprovalStatus(gomock.Any(), uint(1), dto).Return(nil)

	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/additional-repair/1/approval", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCustomerApproval_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUC := mocks.NewMockIAdditionalRepairUseCase(ctrl)
	h := handler.NewAdditionalRepairHandler(mockUC)
	r := setupADRRouter(h)

	dto := entities.AdditionalRepairStatusDTO{ApprovalStatus: "DENIED"}
	mockUC.EXPECT().CustomerApprovalStatus(gomock.Any(), uint(1), dto).Return(errors.New("fail"))

	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/additional-repair/1/approval", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

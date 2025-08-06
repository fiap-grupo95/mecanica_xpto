package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/infrastructure/http/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetVehicles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate1 := valueobject.ParsePlate("ABC1234")
		plate2 := valueobject.ParsePlate("XYZ5678")

		expectedVehicles := []entities.Vehicle{
			{ID: 1, Brand: "Toyota", Model: "Corolla", Year: "2020", Plate: plate1},
			{ID: 2, Brand: "Honda", Model: "Civic", Year: "2021", Plate: plate2},
		}

		mockService.On("GetAllVehicles").Return(expectedVehicles, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetVehicles(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entities.Vehicle
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicles, response)
		mockService.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("GetAllVehicles").Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetVehicles(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetVehiclesByCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := valueobject.ParsePlate("ABC1234")
		expectedVehicles := []entities.Vehicle{
			{ID: 1, Brand: "Toyota", Model: "Corolla", Year: "2020", Plate: plate},
		}

		mockService.On("GetVehiclesByCustomerID", uint(1)).Return(expectedVehicles, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "customerID", Value: "1"}}

		handler.GetVehiclesByCustomerID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entities.Vehicle
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicles, response)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid customer id", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "customerID", Value: "invalid"}}

		handler.GetVehiclesByCustomerID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := valueobject.ParsePlate("ABC1234")
		vehicle := entities.Vehicle{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: plate,
		}

		expectedVehicle := vehicle
		expectedVehicle.ID = 1

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return(&expectedVehicle, nil)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.Vehicle
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle, response)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := valueobject.ParsePlate("ABC1234")
		vehicle := entities.Vehicle{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2021",
			Plate: plate,
		}

		expectedVehicle := vehicle
		expectedVehicle.ID = 1

		mockService.On("UpdateVehicle", mock.AnythingOfType("entities.Vehicle")).Return(&expectedVehicle, nil)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.Vehicle
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle, response)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeleteVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("DeleteVehicle", uint(1)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("DeleteVehicle", uint(1)).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

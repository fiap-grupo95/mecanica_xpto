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

		expected := "Vehicle created successfully"

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return(expected, nil)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected, response)
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

		updates := map[string]interface{}{
			"model": "Corolla",
			"year":  "2021",
		}

		mockService.On("UpdateVehiclePartial", uint(1), updates).Return("Vehicle updated successfully", nil)

		updatesJSON, _ := json.Marshal(updates)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(updatesJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Vehicle updated successfully", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		updates := map[string]interface{}{
			"model": "Corolla",
		}

		updatesJSON, _ := json.Marshal(updates)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(updatesJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		updates := map[string]interface{}{
			"model": "Corolla",
		}

		mockService.On("UpdateVehiclePartial", uint(1), updates).Return("", errors.New("database error"))

		updatesJSON, _ := json.Marshal(updates)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(updatesJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
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

func TestGetVehiclesByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := valueobject.ParsePlate("ABC1234")
		expectedVehicle := &entities.Vehicle{
			ID:    1,
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: plate,
		}

		mockService.On("GetVehicleByID", uint(1)).Return(expectedVehicle, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetVehicleByID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.Vehicle
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle.ID, response.ID)
		assert.Equal(t, expectedVehicle.Brand, response.Brand)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetVehicleByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("GetVehicleByID", uint(1)).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.GetVehicleByID(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetVehiclesByPlate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := "ABC1234"
		expectedVehicle := &entities.Vehicle{
			ID:    1,
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: valueobject.ParsePlate(plate),
		}

		mockService.On("GetVehicleByPlate", plate).Return(expectedVehicle, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: plate}}

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.Vehicle
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle.ID, response.ID)
		assert.Equal(t, string(expectedVehicle.Plate), string(response.Plate))
		mockService.AssertExpectations(t)
	})

	t.Run("empty plate", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: ""}}

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := "ABC1234"
		mockService.On("GetVehicleByPlate", plate).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: plate}}

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

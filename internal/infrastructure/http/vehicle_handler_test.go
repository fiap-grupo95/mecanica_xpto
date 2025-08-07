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

		// Testing both plate formats
		plate1 := valueobject.ParsePlate("ABC1D23") // Mercosul format
		plate2 := valueobject.ParsePlate("XYZ5678") // Old format

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

	t.Run("database error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("GetAllVehicles").Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetVehicles(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "database error")
		mockService.AssertExpectations(t)
	})
}

func TestGetVehicleByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := valueobject.ParsePlate("ABC1D23")
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
		assert.Equal(t, string(expectedVehicle.Plate), string(response.Plate))
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id format", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.GetVehicleByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid vehicle ID")
	})
}

func TestGetVehicleByPlate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success with mercosul format", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := "ABC1D23"
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
		assert.Equal(t, string(expectedVehicle.Plate), string(response.Plate))
		mockService.AssertExpectations(t)
	})

	t.Run("success with old format", func(t *testing.T) {
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
		assert.Equal(t, string(expectedVehicle.Plate), string(response.Plate))
	})

	t.Run("invalid plate format", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: "INVALID"}}

		mockService.On("GetVehicleByPlate", "INVALID").Return(nil, errors.New("invalid plate format"))

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := "ABC1D23"
		mockService.On("GetVehicleByPlate", plate).Return(nil, errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: plate}}

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success with mercosul format", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		vehicle := entities.Vehicle{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: valueobject.ParsePlate("ABC1D23"),
			Customer: entities.Customer{
				ID: 1,
			},
		}

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return("success", nil)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("success with old format", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		vehicle := entities.Vehicle{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: valueobject.ParsePlate("ABC1234"),
			Customer: entities.Customer{
				ID: 1,
			},
		}

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return("success", nil)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("invalid plate format", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		vehicle := entities.Vehicle{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: valueobject.ParsePlate("INVALID"),
		}

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return(nil, errors.New("invalid plate format"))

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		vehicle := entities.Vehicle{
			Brand: "Toyota",
			// Missing other required fields
		}

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return(nil, errors.New("missing required fields"))

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateVehicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		updates := map[string]interface{}{
			"model": "Corolla Updated",
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

	t.Run("vehicle not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		updates := map[string]interface{}{
			"model": "Corolla",
		}

		mockService.On("UpdateVehiclePartial", uint(999), updates).Return("", errors.New("vehicle not found"))

		updatesJSON, _ := json.Marshal(updates)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(updatesJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
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

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("DeleteVehicle", uint(999)).Return(errors.New("not found"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

package http

import (
	"bytes"
	"encoding/json"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/internal/infrastructure/http/mocks"
	"mecanica_xpto/pkg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Função auxiliar para decodificar o corpo da resposta em uma estrutura ErrorResponse
func getErrorResponse(body []byte) pkg.ErrorResponse {
	var response pkg.ErrorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		panic("Erro ao decodificar resposta de erro: " + err.Error())
	}
	return response
}

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

	t.Run("internal error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("GetAllVehicles").Return(nil, pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", nil, http.StatusInternalServerError))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetVehicles(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INTERNAL_ERROR", response.Code)
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
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_VEHICLE_ID", response.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("GetVehicleByID", uint(999)).Return(nil, usecase.ErrVehicleNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.GetVehicleByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "VEHICLE_NOT_FOUND", response.Code)
		mockService.AssertExpectations(t)
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

		mockService.On("GetVehicleByPlate", "INVALID").Return(nil, usecase.ErrInvalidPlateFormat)

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_PLATE_FORMAT", response.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		plate := "ABC1D23"
		mockService.On("GetVehicleByPlate", plate).Return(nil, usecase.ErrVehicleNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "plate", Value: plate}}

		handler.GetVehicleByPlate(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "VEHICLE_NOT_FOUND", response.Code)
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
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "EMPTY_PLATE", response.Code)
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

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return("", usecase.ErrInvalidPlateFormat)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_PLATE_FORMAT", response.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("vehicle already exists", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		vehicle := entities.Vehicle{
			Brand: "Toyota",
			Model: "Corolla",
			Year:  "2020",
			Plate: valueobject.ParsePlate("ABC1234"),
		}

		mockService.On("CreateVehicle", mock.AnythingOfType("entities.Vehicle")).Return("", usecase.ErrVehicleAlreadyExists)

		vehicleJSON, _ := json.Marshal(vehicle)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(vehicleJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "VEHICLE_EXISTS", response.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		invalidJSON := []byte(`{invalid json`)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_INPUT", response.Code)
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
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_VEHICLE_ID", response.Code)
	})

	t.Run("vehicle not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		updates := map[string]interface{}{
			"model": "Corolla",
		}

		mockService.On("UpdateVehiclePartial", uint(999), updates).Return("", usecase.ErrVehicleNotFound)

		updatesJSON, _ := json.Marshal(updates)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(updatesJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "VEHICLE_NOT_FOUND", response.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid plate format in update", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		updates := map[string]interface{}{
			"plate": "INVALID",
		}

		mockService.On("UpdateVehiclePartial", uint(1), updates).Return("", usecase.ErrInvalidPlateFormat)

		updatesJSON, _ := json.Marshal(updates)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest("PATCH", "/", bytes.NewBuffer(updatesJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_PLATE_FORMAT", response.Code)
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
		w.Code = http.StatusNoContent // Expecting no content response
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INVALID_VEHICLE_ID", response.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("DeleteVehicle", uint(999)).Return(usecase.ErrVehicleNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "VEHICLE_NOT_FOUND", response.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("internal error", func(t *testing.T) {
		mockService := new(mocks.MockVehicleService)
		handler := NewVehicleHandler(mockService)

		mockService.On("DeleteVehicle", uint(1)).Return(pkg.NewDomainError("INTERNAL_ERROR", "An internal error occurred", nil, http.StatusInternalServerError))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteVehicle(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		response := getErrorResponse(w.Body.Bytes())
		assert.Equal(t, "INTERNAL_ERROR", response.Code)
		mockService.AssertExpectations(t)
	})
}

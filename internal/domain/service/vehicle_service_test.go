package service

import (
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/service/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetVehicles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockVehicles := []dto.VehicleDTO{
			{
				ID:         1,
				Plate:      "ABC1234",
				Model:      "Civic",
				Brand:      "Honda",
				Year:       "2020",
				CustomerID: 1,
				CreatedAt:  time.Now(),
			},
			{
				ID:         2,
				Plate:      "XYZ5678",
				Model:      "Corolla",
				Brand:      "Toyota",
				Year:       "2021",
				CustomerID: 1,
				CreatedAt:  time.Now(),
			},
		}

		mockRepo.On("FindAll").Return(mockVehicles, nil)

		vehicles, err := service.GetAllVehicles()

		assert.NoError(t, err)
		assert.Len(t, vehicles, 2)
		assert.Equal(t, mockVehicles[0].ID, vehicles[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindAll").Return(nil, errors.New("database error"))

		vehicles, err := service.GetAllVehicles()

		assert.Error(t, err)
		assert.Nil(t, vehicles)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetVehicleByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockVehicleDTO := &dto.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CustomerID: 1,
			CreatedAt:  time.Now(),
		}

		mockRepo.On("FindByID", uint(1)).Return(mockVehicleDTO, nil)

		vehicle, err := service.GetVehicleByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, vehicle)
		assert.Equal(t, mockVehicleDTO.ID, vehicle.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

		vehicle, err := service.GetVehicleByID(1)

		assert.Error(t, err)
		assert.Nil(t, vehicle)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetVehicleByPlate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		plate := valueobject.ParsePlate("ABC1234")
		mockVehicleDTO := &dto.VehicleDTO{
			ID:         1,
			Plate:      string(plate),
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CustomerID: 1,
			CreatedAt:  time.Now(),
		}

		mockRepo.On("FindByPlate", plate).Return(mockVehicleDTO, nil)

		vehicle, err := service.GetVehicleByPlate("ABC1234")

		assert.NoError(t, err)
		assert.NotNil(t, vehicle)
		assert.Equal(t, string(plate), string(vehicle.Plate))
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		plate := valueobject.ParsePlate("XYZ9999")
		mockRepo.On("FindByPlate", plate).Return(nil, errors.New("not found"))

		vehicle, err := service.GetVehicleByPlate("XYZ9999")

		assert.Error(t, err)
		assert.Nil(t, vehicle)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetVehiclesByCustomerID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockVehicles := []dto.VehicleDTO{
			{
				ID:         1,
				Plate:      "ABC1234",
				Model:      "Civic",
				Brand:      "Honda",
				Year:       "2020",
				CustomerID: 1,
				CreatedAt:  time.Now(),
			},
		}

		mockRepo.On("FindByCustomerID", uint(1)).Return(mockVehicles, nil)

		vehicles, err := service.GetVehiclesByCustomerID(1)

		assert.NoError(t, err)
		assert.NotNil(t, vehicles)
		assert.Len(t, vehicles, 1)
		assert.Equal(t, mockVehicles[0].ID, vehicles[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByCustomerID", uint(999)).Return([]dto.VehicleDTO{}, nil)

		vehicles, err := service.GetVehiclesByCustomerID(999)

		assert.NoError(t, err)
		assert.Empty(t, vehicles)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateVehicle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
			Customer: entities.Customer{
				ID: 1,
			},
		}

		mockRepo.On("Create", vehicle).Return(nil)

		result, err := service.CreateVehicle(vehicle)

		assert.NoError(t, err)
		assert.Equal(t, "Vehicle created successfully", result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
			Customer: entities.Customer{
				ID: 1,
			},
		}

		mockRepo.On("Create", vehicle).Return(errors.New("database error"))

		result, err := service.CreateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, "error creating a new vehicle", result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateVehicle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			ID:    1,
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic Updated",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("Update", vehicle).Return(nil)

		result, err := service.UpdateVehicle(vehicle)

		assert.NoError(t, err)
		assert.Equal(t, "Vehicle updated successfully", result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			ID:    1,
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("Update", vehicle).Return(errors.New("database error"))

		result, err := service.UpdateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, "error updating a new vehicle", result)
		mockRepo.AssertExpectations(t)
	})
}

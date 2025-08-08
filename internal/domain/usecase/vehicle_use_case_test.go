package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/usecase/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrorDatabase = errors.New("database error")
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

		mockRepo.On("FindAll").Return(nil, ErrorDatabase)

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

		mockRepo.On("FindByPlate", vehicle.Plate).Return(nil, nil) // No existing vehicle with this plate
		mockRepo.On("Create", vehicle).Return(nil)

		result, err := service.CreateVehicle(vehicle)

		assert.NoError(t, err)
		assert.Equal(t, MessageVehicleCreatedSuccessfully, result)
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

		mockRepo.On("FindByPlate", vehicle.Plate).Return(nil, nil) // No existing vehicle with this plate
		mockRepo.On("Create", vehicle).Return(ErrorDatabase)

		result, err := service.CreateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, MessageErrorCreatingVehicle, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid plate format", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			Plate: valueobject.ParsePlate("ABC123"), // Placa inválida
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		result, err := service.CreateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPlateFormat, err)
		assert.Equal(t, MessageInvalidPlateFormat, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("vehicle already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		existingVehicle := &dto.VehicleDTO{
			ID:    1,
			Plate: "ABC1234",
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("FindByPlate", vehicle.Plate).Return(existingVehicle, nil)

		result, err := service.CreateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, ErrVehicleAlreadyExists, err)
		assert.Equal(t, MessageVehicleAlreadyExists, result)
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

		mockGet := &dto.VehicleDTO{
			ID:    1,
			Plate: "ABC1234",
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("FindByID", vehicle.ID).Return(mockGet, nil)
		mockRepo.On("Update", vehicle).Return(nil)

		result, err := service.UpdateVehicle(vehicle)

		assert.NoError(t, err)
		assert.Equal(t, MessageVehicleUpdatedSuccessfully, result)
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

		mockGet := &dto.VehicleDTO{
			ID:    1,
			Plate: "ABC1234",
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("FindByID", vehicle.ID).Return(mockGet, nil)
		mockRepo.On("Update", vehicle).Return(ErrorDatabase)

		result, err := service.UpdateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, MessageErrorUpdatingVehicle, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("vehicle not found", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			ID:    1,
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("FindByID", vehicle.ID).Return(nil, nil)

		result, err := service.UpdateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, ErrVehicleNotFound, err)
		assert.Equal(t, MessageVehicleNotFound, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid plate format", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		vehicle := entities.Vehicle{
			ID:    1,
			Plate: valueobject.ParsePlate("ABC123"), // Placa inválida
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		result, err := service.UpdateVehicle(vehicle)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPlateFormat, err)
		assert.Equal(t, MessageInvalidPlateFormat, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateVehiclePartial(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		existingVehicle := &dto.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CustomerID: 1,
		}

		updates := map[string]interface{}{
			"model": "Civic Updated",
			"year":  "2021",
		}

		mockRepo.On("FindByID", uint(1)).Return(existingVehicle, nil)
		mockRepo.On("Update", mock.Anything).Return(nil)

		result, err := service.UpdateVehiclePartial(1, updates)

		assert.NoError(t, err)
		assert.Equal(t, MessageVehicleUpdatedSuccessfully, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("vehicle not found", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		updates := map[string]interface{}{
			"model": "Civic Updated",
		}

		mockRepo.On("FindByID", uint(1)).Return(nil, nil)

		result, err := service.UpdateVehiclePartial(1, updates)

		assert.Error(t, err)
		assert.Equal(t, ErrVehicleNotFound, err)
		assert.Equal(t, MessageVehicleNotFound, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid plate format", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		existingVehicle := &dto.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CustomerID: 1,
		}

		updates := map[string]interface{}{
			"plate": "ABC123", // Placa inválida
		}

		mockRepo.On("FindByID", uint(1)).Return(existingVehicle, nil)

		result, err := service.UpdateVehiclePartial(1, updates)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidPlateFormat, err)
		assert.Equal(t, MessageInvalidPlateFormat, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteVehicle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockVehicle := &dto.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CustomerID: 1,
		}

		mockRepo.On("FindByID", uint(1)).Return(mockVehicle, nil)
		mockRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteVehicle(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		err := service.DeleteVehicle(0)

		assert.Error(t, err)
		assert.Equal(t, ErrInvalidID, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("vehicle not found", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByID", uint(1)).Return(nil, nil)

		err := service.DeleteVehicle(1)

		assert.Error(t, err)
		assert.Equal(t, ErrVehicleNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockRepo := new(mocks.MockVehicleRepository)
		service := NewVehicleService(mockRepo)

		mockVehicle := &dto.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CustomerID: 1,
		}

		mockRepo.On("FindByID", uint(1)).Return(mockVehicle, nil)
		mockRepo.On("Delete", uint(1)).Return(ErrorDatabase)

		err := service.DeleteVehicle(1)

		assert.Error(t, err)
		assert.Equal(t, ErrorDatabase, err)
		mockRepo.AssertExpectations(t)
	})
}

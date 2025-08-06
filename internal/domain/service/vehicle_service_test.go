package service

import (
	"errors"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/repository"
	"mecanica_xpto/internal/domain/service/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVehicleService_GetAllVehicles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockVehicles := []repository.VehicleDTO{
			{
				ID:         1,
				Plate:      "ABC1234",
				Model:      "Civic",
				Brand:      "Honda",
				Year:       "2020",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				CustomerID: 1,
			},
			{
				ID:         2,
				Plate:      "XYZ5678",
				Model:      "Corolla",
				Brand:      "Toyota",
				Year:       "2021",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				CustomerID: 1,
			},
		}

		mockRepo.On("FindAll").Return(mockVehicles, nil).Once()

		vehicles, err := service.GetAllVehicles()

		assert.NoError(t, err)
		assert.Len(t, vehicles, 2)
		assert.Equal(t, uint(1), vehicles[0].ID)
		assert.Equal(t, "ABC1234", vehicles[0].Plate.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		var emptySlice []repository.VehicleDTO
		mockRepo.On("FindAll").Return(emptySlice, errors.New("database error")).Once()

		vehicles, err := service.GetAllVehicles()

		assert.Error(t, err)
		assert.Nil(t, vehicles)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehicleByID(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockVehicle := &repository.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			CustomerID: 1,
		}

		mockRepo.On("FindByID", uint(1)).Return(mockVehicle, nil)

		vehicle, err := service.GetVehicleByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, vehicle)
		assert.Equal(t, uint(1), vehicle.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

		vehicle, err := service.GetVehicleByID(999)

		assert.Error(t, err)
		assert.Nil(t, vehicle)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehicleByPlate(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		plate := valueobject.ParsePlate("ABC1234")
		mockVehicle := &repository.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			CustomerID: 1,
		}

		mockRepo.On("FindByPlate", plate).Return(mockVehicle, nil)

		vehicle, err := service.GetVehicleByPlate("ABC1234")

		assert.NoError(t, err)
		assert.NotNil(t, vehicle)
		assert.Equal(t, "ABC1234", vehicle.Plate.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		plate := valueobject.ParsePlate("XYZ9999")
		mockRepo.On("FindByPlate", plate).Return(nil, nil)

		vehicle, err := service.GetVehicleByPlate("XYZ9999")

		assert.NoError(t, err)
		assert.Nil(t, vehicle)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_CreateVehicle(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		vehicle := entities.Vehicle{
			Plate:     valueobject.ParsePlate("ABC1234"),
			Model:     "Civic",
			Brand:     "Honda",
			Year:      "2020",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockVehicleDTO := &repository.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic",
			Brand:      "Honda",
			Year:       "2020",
			CreatedAt:  vehicle.CreatedAt,
			UpdatedAt:  vehicle.UpdatedAt,
			CustomerID: 1,
		}

		mockRepo.On("Create", vehicle).Return(mockVehicleDTO, nil)

		created, err := service.CreateVehicle(vehicle)

		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, uint(1), created.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation error", func(t *testing.T) {
		vehicle := entities.Vehicle{
			Plate: valueobject.ParsePlate("ABC1234"),
			Model: "Civic",
			Brand: "Honda",
			Year:  "2020",
		}

		mockRepo.On("Create", vehicle).Return(nil, errors.New("creation error"))

		created, err := service.CreateVehicle(vehicle)

		assert.Error(t, err)
		assert.Nil(t, created)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_UpdateVehicle(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		vehicle := entities.Vehicle{
			ID:        1,
			Plate:     valueobject.ParsePlate("ABC1234"),
			Model:     "Civic Updated",
			Brand:     "Honda",
			Year:      "2020",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockVehicleDTO := &repository.VehicleDTO{
			ID:         1,
			Plate:      "ABC1234",
			Model:      "Civic Updated",
			Brand:      "Honda",
			Year:       "2020",
			CreatedAt:  vehicle.CreatedAt,
			UpdatedAt:  vehicle.UpdatedAt,
			CustomerID: 1,
		}

		mockRepo.On("Update", vehicle).Return(mockVehicleDTO, nil)

		updated, err := service.UpdateVehicle(vehicle)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "Civic Updated", updated.Model)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		vehicle := entities.Vehicle{
			ID:    999,
			Plate: valueobject.ParsePlate("ABC1234"),
		}

		mockRepo.On("Update", vehicle).Return(nil, errors.New("update error"))

		updated, err := service.UpdateVehicle(vehicle)

		assert.Error(t, err)
		assert.Nil(t, updated)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_DeleteVehicle(t *testing.T) {
	mockRepo := mocks.NewMockVehicleRepository()
	service := NewVehicleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", uint(1)).Return(nil)

		err := service.DeleteVehicle(1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete error", func(t *testing.T) {
		mockRepo.On("Delete", uint(999)).Return(errors.New("delete error"))

		err := service.DeleteVehicle(999)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestVehicleService_GetVehiclesByCustomerID(t *testing.T) {
	t.Run("success with vehicles", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockVehicles := []repository.VehicleDTO{
			{
				ID:         1,
				CustomerID: 1,
				Plate:      "ABC1234",
				Model:      "Civic",
				Brand:      "Honda",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				ID:         2,
				CustomerID: 1,
				Plate:      "XYZ5678",
				Model:      "Corolla",
				Brand:      "Toyota",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		mockRepo.On("FindByCustomerID", uint(1)).Return(mockVehicles, nil).Once()

		vehicles, err := service.GetVehiclesByCustomerID(1)

		assert.NoError(t, err)
		assert.Len(t, vehicles, 2)
		assert.Equal(t, uint(1), vehicles[0].ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success with no vehicles", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		mockRepo.On("FindByCustomerID", uint(999)).Return([]repository.VehicleDTO{}, nil).Once()

		vehicles, err := service.GetVehiclesByCustomerID(999)

		assert.NoError(t, err)
		assert.Empty(t, vehicles)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := mocks.NewMockVehicleRepository()
		service := NewVehicleService(mockRepo)

		var emptySlice []repository.VehicleDTO
		mockRepo.On("FindByCustomerID", uint(1)).Return(emptySlice, errors.New("database error")).Once()

		vehicles, err := service.GetVehiclesByCustomerID(1)

		assert.Error(t, err)
		assert.Nil(t, vehicles)
		mockRepo.AssertExpectations(t)
	})
}

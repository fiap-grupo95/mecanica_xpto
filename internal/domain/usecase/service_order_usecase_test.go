package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/usecase/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Vehicle Repository - "github.com/stretchr/testify/mock"
type MockVehicleRepository struct {
	mock.Mock
}

func (m *MockVehicleRepository) FindAll() ([]dto.VehicleDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) FindByID(id uint) (*dto.VehicleDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) FindByPlate(plate valueobject.Plate) (*dto.VehicleDTO, error) {
	args := m.Called(plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) FindByCustomerID(customerID uint) ([]dto.VehicleDTO, error) {
	args := m.Called(customerID)
	return args.Get(0).([]dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) Create(vehicle entities.Vehicle) error {
	args := m.Called(vehicle)
	return args.Error(0)
}

func (m *MockVehicleRepository) Update(vehicle entities.Vehicle) error {
	args := m.Called(vehicle)
	return args.Error(0)
}

func (m *MockVehicleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVehicleRepository) GetByID(id uint) (*entities.Vehicle, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Vehicle), args.Error(1)
}

// Mock Customer Repository - "github.com/stretchr/testify/mock"
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(customer *dto.CustomerDTO) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetByID(id uint) (*dto.CustomerDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CustomerDTO), args.Error(1)
}

func (m *MockCustomerRepository) Update(customer *dto.CustomerDTO) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetByDocument(CpfCnpj string) (*dto.CustomerDTO, error) {
	args := m.Called(CpfCnpj)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CustomerDTO), args.Error(1)
}

func (m *MockCustomerRepository) List() ([]dto.CustomerDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.CustomerDTO), args.Error(1)
}

func TestCreateServiceOrder(t *testing.T) {
	// Setup mocks
	mockServiceOrderRepo := new(mocks.MockServiceOrderRepository)
	mockVehicleRepo := new(MockVehicleRepository)
	mockCustomerRepo := new(MockCustomerRepository)

	// Create use case instance
	useCase := NewServiceOrderUseCase(mockServiceOrderRepo, mockVehicleRepo, mockCustomerRepo)

	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success - Valid service order creation",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 2,
				VehicleID:  2,
			},
			setupMocks: func() {
				mockVehicleRepo.On("FindByID", uint(2)).Return(&dto.VehicleDTO{ID: 2}, nil)
				mockCustomerRepo.On("GetByID", uint(2)).Return(&dto.CustomerDTO{ID: 2}, nil)
				mockServiceOrderRepo.On("Create", mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid vehicle ID",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 2,
				VehicleID:  999,
			},
			setupMocks: func() {
				mockVehicleRepo.On("FindByID", uint(999)).Return(nil, errors.New("vehicle not found"))
			},
			expectedError: errors.New("vehicle not found"),
		},
		{
			name: "Error - Invalid customer ID",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 999,
				VehicleID:  2,
			},
			setupMocks: func() {
				mockVehicleRepo.On("FindByID", uint(2)).Return(&dto.VehicleDTO{ID: 2}, nil)
				mockCustomerRepo.On("GetByID", uint(999)).Return(nil, errors.New("customer not found"))
			},
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks for this test case
			tt.setupMocks()

			// Execute the test
			err := useCase.CreateServiceOrder(tt.serviceOrder)

			// Verify the results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify that all expected mock calls were made
			mockVehicleRepo.AssertExpectations(t)
			mockCustomerRepo.AssertExpectations(t)
			mockServiceOrderRepo.AssertExpectations(t)
		})
	}
}

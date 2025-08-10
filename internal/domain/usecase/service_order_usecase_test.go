package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Vehicle Repository - "github.com/stretchr/testify/mock"
type MockVehicleRepository struct {
	mock.Mock
}

func (m *MockVehicleRepository) FindAll(ctx context.Context) ([]dto.VehicleDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) FindByID(ctx context.Context, id uint) (*dto.VehicleDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) FindByPlate(ctx context.Context, plate valueobject.Plate) (*dto.VehicleDTO, error) {
	args := m.Called(ctx, plate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) FindByCustomerID(ctx context.Context, customerID uint) ([]dto.VehicleDTO, error) {
	args := m.Called(ctx, customerID)
	return args.Get(0).([]dto.VehicleDTO), args.Error(1)
}

func (m *MockVehicleRepository) Create(ctx context.Context, vehicle *entities.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

func (m *MockVehicleRepository) Update(ctx context.Context, vehicle *entities.Vehicle) error {
	args := m.Called(ctx, vehicle)
	return args.Error(0)
}

// Mock Customer Repository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetByID(id uint) (*dto.CustomerDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CustomerDTO), args.Error(1)
}

// Mock Service Order Repository
type MockServiceOrderRepository struct {
	mock.Mock
}

func (m *MockServiceOrderRepository) Create(serviceOrder *entities.ServiceOrder) error {
	args := m.Called(serviceOrder)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockServiceOrderRepository) List() ([]dto.ServiceOrderDTO, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ServiceOrderDTO), args.Error(1)
}

func (m *MockServiceOrderRepository) GetStatus(status valueobject.ServiceOrderStatus) (*dto.ServiceOrderStatusDTO, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ServiceOrderStatusDTO), args.Error(1)
}

func (m *MockServiceOrderRepository) Update(serviceOrder *entities.ServiceOrder) error {
	args := m.Called(serviceOrder)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(0)
}

func (m *MockServiceOrderRepository) GetByID(id uint) (*dto.ServiceOrderDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ServiceOrderDTO), args.Error(1)
}

// Mock Parts Supply Repository
type MockPartsSupplyRepository struct {
	mock.Mock
}

func (m *MockPartsSupplyRepository) GetByID(ctx context.Context, id uint) (*entities.PartsSupply, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PartsSupply), args.Error(1)
}

func (m *MockPartsSupplyRepository) Update(ctx context.Context, ps *entities.PartsSupply) error {
	args := m.Called(ctx, ps)
	return args.Error(0)
}

// Mock Service Repository
type ServiceRepoMock struct {
	mock.Mock
}

func (m *ServiceRepoMock) GetByID(ctx context.Context, id uint) (entities.Service, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return entities.Service{}, args.Error(1)
	}
	return args.Get(0).(entities.Service), args.Error(1)
}

func (m *ServiceRepoMock) Create(ctx context.Context, service *entities.Service) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *ServiceRepoMock) Update(ctx context.Context, service *entities.Service) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *ServiceRepoMock) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ServiceRepoMock) List(ctx context.Context) ([]entities.Service, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entities.Service), args.Error(1)
}

func TestCreateServiceOrder(t *testing.T) {
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceOrderRepo := new(MockServiceOrderRepository)
	serviceRepo := new(ServiceRepoMock)
	partsSupplyRepo := new(MockPartsSupplyRepository)

	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success - Valid service order creation",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 1,
				VehicleID:  1,
			},
			setupMocks: func() {
				vehicleRepo.On("FindByID", mock.Anything, uint(1)).Return(&dto.VehicleDTO{ID: 1}, nil)
				customerRepo.On("GetByID", uint(1)).Return(&dto.CustomerDTO{ID: 1}, nil)
				serviceOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error - Vehicle not found",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 1,
				VehicleID:  1,
			},
			setupMocks: func() {
				vehicleRepo.On("FindByID", mock.Anything, uint(1)).Return(nil, errors.New("vehicle not found"))
			},
			expectedError: errors.New("vehicle not found"),
		},
		{
			name: "Error - Customer not found",
			serviceOrder: entities.ServiceOrder{
				CustomerID: 1,
				VehicleID:  1,
			},
			setupMocks: func() {
				vehicleRepo.On("FindByID", mock.Anything, uint(1)).Return(&dto.VehicleDTO{ID: 1}, nil)
				customerRepo.On("GetByID", uint(1)).Return(nil, errors.New("customer not found"))
			},
			expectedError: errors.New("customer not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := useCase.CreateServiceOrder(context.Background(), tt.serviceOrder)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateServiceOrder(t *testing.T) {
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceOrderRepo := new(MockServiceOrderRepository)
	serviceRepo := new(ServiceRepoMock)
	partsSupplyRepo := new(MockPartsSupplyRepository)

	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	setupMocks := func() {
		serviceOrderRepo.On("GetByID", uint(1)).Return(&dto.ServiceOrderDTO{
			ID: 1,
			ServiceOrderStatus: dto.ServiceOrderStatusDTO{
				ID:          1,
				Description: string(valueobject.StatusRecebida),
			},
		}, nil)
		serviceRepo.On("GetByID", uint(1)).Return(&dto.ServiceDTO{ID: 1}, nil)
		serviceOrderRepo.On("Update", mock.Anything, mock.AnythingOfType("*entities.ServiceOrder")).Return(nil)
	}
	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		flow          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Success - Update to EmDiagnostico",
			serviceOrder: entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEmDiagnostico,
				Services: []entities.Service{
					{ID: 1},
				},
			},
			flow:          DIAGNOSIS,
			setupMocks:    setupMocks,
			expectedError: nil,
		},
		{
			name: "Error - Invalid Status Transition",
			serviceOrder: entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
			},
			flow: DIAGNOSIS,
			setupMocks: func() {
				serviceOrderRepo.On("GetByID", uint(1)).Return(&dto.ServiceOrderDTO{
					ID: 1,
					ServiceOrderStatus: dto.ServiceOrderStatusDTO{
						ID:          1,
						Description: string(valueobject.StatusRecebida),
					},
				}, nil)
			},
			expectedError: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := useCase.UpdateServiceOrder(context.Background(), tt.serviceOrder, tt.flow)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEstimate(t *testing.T) {
	tests := []struct {
		name            string
		request         *entities.ServiceOrder
		serviceOrderDTO *dto.ServiceOrderDTO
		partsSupplies   []entities.PartsSupply
		setupMocks      func(*MockPartsSupplyRepository)
		expectedError   error
	}{
		{
			name: "Success - Approve Estimate",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusAprovada,
				PartsSupplies: []entities.PartsSupply{
					{ID: 1, QuantityReserve: 2},
				},
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					ID:          1,
					Description: string(valueobject.StatusAguardandoAprovacao),
				},
			},
			setupMocks: func(repo *MockPartsSupplyRepository) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(&dto.PartsSupplyDTO{
					ID:              1,
					QuantityTotal:   5,
					QuantityReserve: 2,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*entities.PartsSupply")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Success - Reject Estimate",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusRejeitada,
				PartsSupplies: []entities.PartsSupply{
					{ID: 1, QuantityReserve: 2},
				},
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					ID:          1,
					Description: string(valueobject.StatusAguardandoAprovacao),
				},
			},
			setupMocks: func(repo *MockPartsSupplyRepository) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(&dto.PartsSupplyDTO{
					ID:              1,
					QuantityTotal:   5,
					QuantityReserve: 2,
				}, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*entities.PartsSupply")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			partsSupplyRepo := new(MockPartsSupplyRepository)
			if tt.setupMocks != nil {
				tt.setupMocks(partsSupplyRepo)
			}

			update := &entities.ServiceOrder{}
			result, err := ValidateEstimate(context.Background(), tt.request, tt.serviceOrderDTO, update, partsSupplyRepo)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestCalculateEstimate(t *testing.T) {
	tests := []struct {
		name          string
		services      []entities.Service
		partsSupplies []entities.PartsSupply
		expected      float64
	}{
		{
			name: "Calculate with services and parts supplies",
			services: []entities.Service{
				{ID: 1, Price: 100.0},
				{ID: 2, Price: 150.0},
			},
			partsSupplies: []entities.PartsSupply{
				{ID: 1, Price: 50.0, QuantityTotal: 2},
				{ID: 2, Price: 75.0, QuantityReserve: 1},
			},
			expected: 375.0, // (100 + 150) + (50*2 + 75*1)
		},
		{
			name: "Calculate with only services",
			services: []entities.Service{
				{ID: 1, Price: 100.0},
				{ID: 2, Price: 150.0},
			},
			partsSupplies: []entities.PartsSupply{},
			expected:      250.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateEstimate(tt.services, tt.partsSupplies)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateExecution(t *testing.T) {
	tests := []struct {
		name            string
		request         *entities.ServiceOrder
		serviceOrderDTO *dto.ServiceOrderDTO
		expectedError   error
	}{
		{
			name: "Success - Start Execution",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEmExecucao,
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					Description: string(valueobject.StatusAprovada),
				},
			},
			expectedError: nil,
		},
		{
			name: "Success - Finish Execution",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusFinalizada,
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					Description: string(valueobject.StatusEmExecucao),
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Invalid Status Transition",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusFinalizada,
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					Description: string(valueobject.StatusRecebida),
				},
			},
			expectedError: ErrInvalidTransitionStatusToExecution,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &entities.ServiceOrder{}
			result, err := ValidateExecution(context.Background(), tt.request, tt.serviceOrderDTO, update)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.ServiceOrderStatus, result.ServiceOrderStatus)
			}
		})
	}
}

func TestValidateDelivery(t *testing.T) {
	tests := []struct {
		name            string
		request         *entities.ServiceOrder
		serviceOrderDTO *dto.ServiceOrderDTO
		expectedError   error
	}{
		{
			name: "Success - Complete Delivery",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
				Payment: &entities.Payment{
					ID:           1,
					ServiceOrder: entities.ServiceOrder{ID: 1},
					PaymentDate:  time.Now(),
				},
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					Description: string(valueobject.StatusFinalizada),
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - Missing Payment Information",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					Description: string(valueobject.StatusFinalizada),
				},
			},
			expectedError: errors.New("payment information is required for delivery"),
		},
		{
			name: "Error - Invalid Status Transition",
			request: &entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEntregue,
				Payment: &entities.Payment{
					ID:           1,
					ServiceOrder: entities.ServiceOrder{ID: 1},
					PaymentDate:  time.Now(),
				},
			},
			serviceOrderDTO: &dto.ServiceOrderDTO{
				ID: 1,
				ServiceOrderStatus: dto.ServiceOrderStatusDTO{
					Description: string(valueobject.StatusEmDiagnostico),
				},
			},
			expectedError: ErrInvalidTransitionStatusToDelivery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update := &entities.ServiceOrder{}
			result, err := ValidateDelivery(context.Background(), tt.request, tt.serviceOrderDTO, update)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, valueobject.StatusEntregue, result.ServiceOrderStatus)
			}
		})
	}
}

func TestInvalidServiceOrder(t *testing.T) {
	vehicleRepo := new(MockVehicleRepository)
	customerRepo := new(MockCustomerRepository)
	serviceOrderRepo := new(MockServiceOrderRepository)
	serviceRepo := new(ServiceRepoMock)
	partsSupplyRepo := new(MockPartsSupplyRepository)

	useCase := NewServiceOrderUseCase(serviceOrderRepo, vehicleRepo, customerRepo, serviceRepo, partsSupplyRepo)

	tests := []struct {
		name          string
		serviceOrder  entities.ServiceOrder
		flow          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "Error - Service Order Not Found",
			serviceOrder: entities.ServiceOrder{
				ID:                 999,
				ServiceOrderStatus: valueobject.StatusEmDiagnostico,
			},
			flow: DIAGNOSIS,
			setupMocks: func() {
				serviceOrderRepo.On("GetByID", uint(999)).Return(nil, ErrServiceOrderNotFound)
			},
			expectedError: ErrServiceOrderNotFound,
		},
		{
			name: "Error - Invalid Flow Type",
			serviceOrder: entities.ServiceOrder{
				ID:                 1,
				ServiceOrderStatus: valueobject.StatusEmDiagnostico,
			},
			flow: "invalid_flow",
			setupMocks: func() {
				serviceOrderRepo.On("GetByID", uint(1)).Return(&dto.ServiceOrderDTO{
					ID: 1,
					ServiceOrderStatus: dto.ServiceOrderStatusDTO{
						Description: string(valueobject.StatusRecebida),
					},
				}, nil)
			},
			expectedError: nil, // The update will return nil since no valid flow was matched
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := useCase.UpdateServiceOrder(context.Background(), tt.serviceOrder, tt.flow)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

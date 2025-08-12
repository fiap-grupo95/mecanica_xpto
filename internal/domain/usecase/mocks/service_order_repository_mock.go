package mocks

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"

	"github.com/stretchr/testify/mock"
)

// Mock ServiceOrder Repository
type MockServiceOrderRepository struct {
	mock.Mock
}

func (m *MockServiceOrderRepository) Create(serviceOrder *entities.ServiceOrder) error {
	args := m.Called(serviceOrder)
	return args.Error(0)
}

func (m *MockServiceOrderRepository) Update(serviceOrder *entities.ServiceOrder) error {
	args := m.Called(serviceOrder)
	return args.Error(0)
}

func (m *MockServiceOrderRepository) GetByID(id uint) (*dto.ServiceOrderDTO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ServiceOrderDTO), args.Error(1)
}

func (m *MockServiceOrderRepository) List() ([]dto.ServiceOrderDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.ServiceOrderDTO), args.Error(1)
}

func (m *MockServiceOrderRepository) GetStatus(status valueobject.ServiceOrderStatus) (*dto.ServiceOrderStatusDTO, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ServiceOrderStatusDTO), args.Error(1)
}

func (m *MockServiceOrderRepository) GetPartsSupplyServiceOrder(serviceOrderID uint, partsSupplyID uint) (*dto.PartsSupplyServiceOrderDTO, error) {
	args := m.Called(serviceOrderID, partsSupplyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PartsSupplyServiceOrderDTO), args.Error(1)
}

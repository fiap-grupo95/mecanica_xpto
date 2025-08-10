package usecase

import (
	"github.com/rs/zerolog/log"

	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	customerRepo "mecanica_xpto/internal/domain/repository/customers"
	"mecanica_xpto/internal/domain/repository/service_order"
	"mecanica_xpto/internal/domain/repository/vehicles"
)

type IServiceOrderUseCase interface {
	CreateServiceOrder(serviceOrder entities.ServiceOrder) error
	UpdateServiceOrder(serviceOrder entities.ServiceOrder) error
}

type ServiceOrderUseCase struct {
	repo         serviceorder.IServiceOrderRepository
	vehicleRepo  vehicles.VehicleRepositoryInterface
	customerRepo customerRepo.ICustomerRepository
}

func NewServiceOrderUseCase(repo serviceorder.IServiceOrderRepository,
	vehicleRepo vehicles.VehicleRepositoryInterface,
	customerRepo customerRepo.ICustomerRepository) *ServiceOrderUseCase {
	return &ServiceOrderUseCase{
		repo:         repo,
		vehicleRepo:  vehicleRepo,
		customerRepo: customerRepo,
	}
}

// CreateServiceOrder creates a new service order after validating the vehicle and customer.
// It sets the initial status of the service order to "Recebida".
// If the vehicle or customer validation fails, it logs the error and returns it.
func (u *ServiceOrderUseCase) CreateServiceOrder(serviceOrder entities.ServiceOrder) error {
	err := validateVehicle(serviceOrder, u.vehicleRepo)
	if err != nil {
		log.Error().Msgf("Error validating vehicle: %v", err)
		return err
	}

	err = validateCustomer(serviceOrder, u.customerRepo)
	if err != nil {
		log.Error().Msgf("Error validating customer: %v", err)
		return err
	}

	newServiceOrder := entities.ServiceOrder{
		CustomerID:         serviceOrder.CustomerID,
		VehicleID:          serviceOrder.VehicleID,
		ServiceOrderStatus: valueobject.StatusRecebida,
	}

	err = u.repo.Create(&newServiceOrder)
	if err != nil {
		log.Error().Msgf("Error creating service order: %v", err)
		return err
	}
	return nil
}

// UpdateServiceOrder updates an existing service order.
func (u *ServiceOrderUseCase) UpdateServiceOrder(serviceOrder entities.ServiceOrder) error {

	if serviceOrder.ID != 0 {
		return u.repo.Update(&serviceOrder)
	}
	return nil
}

func validateCustomer(serviceOrder entities.ServiceOrder, customerRepo customerRepo.ICustomerRepository) error {
	if serviceOrder.CustomerID != 0 {
		customer, err := customerRepo.GetByID(serviceOrder.CustomerID)
		if err != nil {
			log.Error().Msgf("error finding customer with id %d: %v", serviceOrder.CustomerID, err)
			return err
		}
		if customer == nil {
			return ErrCustomerNotFound
		}
		return nil
	}
	return ErrInvalidCustomerID
}

func validateVehicle(serviceOrder entities.ServiceOrder, vehicleRepo vehicles.VehicleRepositoryInterface) error {
	if serviceOrder.VehicleID != 0 {
		vehicle, err := vehicleRepo.FindByID(serviceOrder.VehicleID)
		if err != nil {
			log.Error().Msgf("error finding vehicle with id %d: %v", serviceOrder.VehicleID, err)
			return err
		}
		if vehicle == nil {
			return ErrVehicleNotFound
		}
		return nil
	}
	return ErrInvalidID
}

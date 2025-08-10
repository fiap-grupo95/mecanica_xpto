package usecase

import (
	"context"
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/repository/parts_supply"

	"github.com/rs/zerolog/log"

	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	customerRepo "mecanica_xpto/internal/domain/repository/customers"
	"mecanica_xpto/internal/domain/repository/service"
	"mecanica_xpto/internal/domain/repository/service_order"
	"mecanica_xpto/internal/domain/repository/vehicles"
)

// operation flow
const (
	DIAGNOSIS = "diagnosis"
	ESTIMATE  = "estimate"
	EXECUTION = "execution"
	DELIVERY  = "delivery"
)

var (
	ErrServiceOrderNotFound               = errors.New("service order not found")
	ErrInvalidTransitionStatusToDiagnosis = errors.New("invalid transition status to diagnosis")
	ErrInvalidStatus                      = errors.New("invalid service order status")
)

type IServiceOrderUseCase interface {
	CreateServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder) error
	UpdateServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder, flow string) error
}

type ServiceOrderUseCase struct {
	repo            serviceorder.IServiceOrderRepository
	vehicleRepo     vehicles.VehicleRepositoryInterface
	customerRepo    customerRepo.ICustomerRepository
	serviceRepo     service.IServiceRepo
	partsSupplyRepo parts_supply.IPartsSupplyRepo
}

func NewServiceOrderUseCase(repo serviceorder.IServiceOrderRepository,
	vehicleRepo vehicles.VehicleRepositoryInterface,
	customerRepo customerRepo.ICustomerRepository,
	serviceRepo service.IServiceRepo,
	partsSupplyRepo parts_supply.IPartsSupplyRepo) *ServiceOrderUseCase {
	return &ServiceOrderUseCase{
		repo:            repo,
		vehicleRepo:     vehicleRepo,
		customerRepo:    customerRepo,
		serviceRepo:     serviceRepo,
		partsSupplyRepo: partsSupplyRepo,
	}
}

// CreateServiceOrder creates a new service order after validating the vehicle and customer.
// It sets the initial status of the service order to "Recebida".
// If the vehicle or customer validation fails, it logs the error and returns it.
func (u *ServiceOrderUseCase) CreateServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder) error {
	err := validateVehicle(ctx, serviceOrder, u.vehicleRepo)
	if err != nil {
		log.Error().Msgf("Error validating vehicle: %v", err)
		return err
	}

	err = validateCustomer(ctx, serviceOrder, u.customerRepo)
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
func (u *ServiceOrderUseCase) UpdateServiceOrder(ctx context.Context, request entities.ServiceOrder, flow string) error {
	var update = &entities.ServiceOrder{}

	if request.ID == 0 {
		return ErrInvalidID
	}

	serviceOrderDto, err := u.repo.GetByID(request.ID)
	if err != nil {
		log.Error().Msgf("Error finding service order with id %v: %v", request, err)
		return err
	}
	if serviceOrderDto == nil {
		log.Error().Msgf("Service order with id %d not found", request.ID)
		return ErrServiceOrderNotFound
	}

	if flow == DIAGNOSIS {
		update, err = ValidateDiagnosis(ctx, &request, serviceOrderDto, update, u.serviceRepo, u.partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error validating diagnosis: %v", err)
		}
	} else if flow == ESTIMATE {
		update, err = ValidateEstimate(ctx, &request, serviceOrderDto, update)
		if err != nil {
			log.Error().Msgf("Error validating estimate: %v", err)
		}
	} else if flow == EXECUTION {
		update, err = ValidateExecution(ctx, &request, serviceOrderDto, update)
		if err != nil {
			log.Error().Msgf("Error validating execution: %v", err)

		}
	} else if flow == DELIVERY {
		update, err = ValidateDelivery(ctx, &request, serviceOrderDto, update)
		if err != nil {
			log.Error().Msgf("Error validating delivery: %v", err)
		}
	}

	return u.repo.Update(update)
}

// ValidateDiagnosis checks if the service order status is valid for diagnosis.
// If the status is "Recebida" or "EmDiagnostico" and the request status is "EmDiagnostico",
// it updates the service order status to "EmDiagnostico".
func ValidateDiagnosis(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder, serviceRepo service.IServiceRepo, partsSupplyRepo parts_supply.IPartsSupplyRepo) (*entities.ServiceOrder, error) {
	newStatus := request.ServiceOrderStatus
	oldStatus := serviceOrderDto.ServiceOrderStatus.ToDomain()

	if !newStatus.IsValid() {
		return nil, ErrInvalidStatus
	}

	if oldStatus == valueobject.StatusRecebida && newStatus == valueobject.StatusCancelada {
		update.ServiceOrderStatus = valueobject.StatusCancelada
		return update, nil
	}

	if (oldStatus == valueobject.StatusRecebida && newStatus == valueobject.StatusEmDiagnostico) || oldStatus == valueobject.StatusEmDiagnostico {
		update.ServiceOrderStatus = valueobject.StatusEmDiagnostico

		if len(request.Services) > 0 {
			update.Services = request.Services

			// Validate if each Service exists
			for _, s := range request.Services {
				err := validateService(ctx, s, serviceRepo)
				if err != nil {
					log.Error().Msgf("Error validating service: %v", err)
					return nil, err
				}
			}
		} else {
			log.Error().Msg("No services provided for diagnosis")
			return nil, errors.New("no services provided for diagnosis")
		}

		if len(request.PartsSupplies) > 0 {
			update.PartsSupplies = request.PartsSupplies

			// Validate if each PartsSupplies are available
			// if all exists and ara available, reserve each PartsSupplies
			for _, ps := range request.PartsSupplies {
				err := validatePartsSupply(ctx, ps, partsSupplyRepo)
				if err != nil {
					log.Error().Msgf("Error validating parts supply: %v", err)
					return nil, err
				}
				if !validatePartsSupply.IsAvailable() {
					log.Error().Msgf("Parts supply with ID %d is not available", ps.ID)
					return nil, errors.New("parts supply not available")
				}
				// Reserve the parts supply
				err = serviceorder.ReservePartsSupply(ps)
				if err != nil {
					log.Error().Msgf("Error reserving parts supply: %v", err)
					return nil, err
				}
			}
		} else {
			log.Error().Msg("No parts supplies provided for diagnosis")
			return nil, errors.New("no parts supplies provided for diagnosis")
		}

		// Calculate the total cost of PartsSupplies and Services and set it to the estimate

		// If OK, set a new status to "AguardandoAprovacao"
	}

	return nil, ErrInvalidTransitionStatusToDiagnosis

}

func ValidateEstimate(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	if !request.ServiceOrderStatus.IsValid() {
		return nil, ErrInvalidStatus
	}
}

func ValidateExecution(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	if !request.ServiceOrderStatus.IsValid() {
		return nil, ErrInvalidStatus
	}
}

func ValidateDelivery(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	if !request.ServiceOrderStatus.IsValid() {
		return nil, ErrInvalidStatus
	}
}

func validateVehicle(ctx context.Context, serviceOrder entities.ServiceOrder, vehicleRepo vehicles.VehicleRepositoryInterface) error {
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

func validateCustomer(ctx context.Context, serviceOrder entities.ServiceOrder, customerRepo customerRepo.ICustomerRepository) error {
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

func validateService(ctx context.Context, s entities.Service, serviceRepo service.IServiceRepo) error {
	// Here you would implement the logic to validate the service
	// For example, check if the service exists in the database
	// and if it is available for the service order.
	// This is a placeholder implementation.
	if s.ID == 0 {
		return ErrInvalidID
	}
	result, err := serviceRepo.GetByID(ctx, s.ID)
	if err != nil {
		log.Error().Msgf("error finding service with id %d: %v", s.ID, err)
		return err
	}
	if result.ID == 0 {
		log.Error().Msgf("service with id %d not found", s.ID)
		return ErrServiceNotFound
	}
	return nil
}

func validatePartsSupply(ctx context.Context, partsSupply entities.PartsSupply, partsSupplyRepo parts_supply.IPartsSupplyRepo) error {
	// Here you would implement the logic to validate the parts supply
	// For example, check if the parts supply exists in the database
	// and if it is available for the service order.
	// This is a placeholder implementation.
	if partsSupply.ID == 0 {
		return ErrInvalidID
	}
	return nil
}

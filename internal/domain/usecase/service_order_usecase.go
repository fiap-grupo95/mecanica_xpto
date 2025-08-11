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
	ErrInvalidTransitionStatusToExecution = errors.New("invalid transition status to execution")
	ErrInvalidTransitionStatusToDelivery  = errors.New("invalid transition status to delivery")
	ErrInvalidTransitionStatusToEstimate  = errors.New("invalid transition status to estimate")
	ErrInvalidStatus                      = errors.New("invalid service order status")
	ErrInsufficientPartsSupply            = errors.New("insufficient parts supply available")
	ErrInvalidFlow                        = errors.New("invalid flow")
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

func NewServiceOrderUseCase(repo serviceorder.IServiceOrderRepository, vehicleRepo vehicles.VehicleRepositoryInterface, customerRepo customerRepo.ICustomerRepository, serviceRepo service.IServiceRepo, partsSupplyRepo parts_supply.IPartsSupplyRepo) *ServiceOrderUseCase {
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
		log.Error().Msgf("Error finding service order with id %v: %v", request.ID, err)
		return err
	}
	if serviceOrderDto == nil {
		log.Error().Msgf("Service order with id %d not found", request.ID)
		return ErrServiceOrderNotFound
	}

	switch flow {
	case DIAGNOSIS:
		update, err = ValidateDiagnosis(ctx, &request, serviceOrderDto, update, u.serviceRepo, u.partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error validating diagnosis: %v", err)
			return err
		}
	case ESTIMATE:
		update, err = ValidateEstimate(ctx, &request, serviceOrderDto, update, u.partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error validating estimate: %v", err)
			return err
		}
	case EXECUTION:
		update, err = ValidateExecution(ctx, &request, serviceOrderDto, update)
		if err != nil {
			log.Error().Msgf("Error validating execution: %v", err)
			return err
		}
	case DELIVERY:
		update, err = ValidateDelivery(ctx, &request, serviceOrderDto, update)
		if err != nil {
			log.Error().Msgf("Error validating delivery: %v", err)
			return err
		}
	default:
		log.Error().Msgf("Invalid flow: %s", flow)
		return ErrInvalidFlow
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
				_, err := getSeviceById(ctx, s, serviceRepo)
				if err != nil {
					log.Error().Msgf("Error getting service by id %v: %v", s.ID, err)
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
				err := validateQttPartsSupply(ctx, ps, partsSupplyRepo)
				if err != nil {
					log.Error().Msgf("Error validating parts supply: %v", err)
					return nil, err
				}
			}

			for _, ps := range request.PartsSupplies {
				// Reserve the parts supply
				err := reservePartsSupply(ctx, ps, partsSupplyRepo)
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
		update.Estimate = CalculateEstimate(update.Services, update.PartsSupplies)

		// If OK, set a new status to "AguardandoAprovacao"
		update.ServiceOrderStatus = valueobject.StatusAguardandoAprovacao
		return update, nil
	}

	return nil, ErrInvalidTransitionStatusToDiagnosis

}

func CalculateEstimate(services []entities.Service, partsSupplies []entities.PartsSupply) float64 {
	totalEstimate := 0.0

	// Calculate services total
	for _, s := range services {
		totalEstimate += s.Price
	}

	// Calculate parts supplies total
	for _, ps := range partsSupplies {
		quantity := ps.QuantityReserve // Use reserve quantity by default
		if ps.QuantityTotal > 0 {      // If total quantity is specified, use that instead
			quantity = ps.QuantityTotal
		}
		totalEstimate += ps.Price * float64(quantity)
	}

	return totalEstimate
}

func ValidateEstimate(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder, partsSupplyRepo parts_supply.IPartsSupplyRepo) (*entities.ServiceOrder, error) {
	oldStatus := serviceOrderDto.ServiceOrderStatus.ToDomain()

	if !request.ServiceOrderStatus.IsValid() {
		return nil, ErrInvalidStatus
	}

	if oldStatus.IsAguardandoAprovacao() && request.ServiceOrderStatus.IsAprovada() {
		update.ServiceOrderStatus = valueobject.StatusAprovada
		// If the status is "Aprovada", we can subtract the total available quantity of PartsSupplies from the quantity reserve
		for _, ps := range request.PartsSupplies {
			err := releaseReservedPartsSupply(ctx, ps, partsSupplyRepo)
			if err != nil {
				log.Error().Msgf("Error releasing reserved parts supply: %v", err)
				return nil, err
			}
		}
		return update, nil
	}
	if oldStatus.IsAguardandoAprovacao() && request.ServiceOrderStatus.IsRejeitada() {
		update.ServiceOrderStatus = valueobject.StatusRejeitada
		// If the status is "Rejeitada", we can reset the PartsSupplies reserve
		for _, ps := range request.PartsSupplies {
			err := unreservePartsSupply(ctx, ps, partsSupplyRepo)
			if err != nil {
				log.Error().Msgf("Error unreserving parts supply: %v", err)
				return nil, err
			}
		}
		return update, nil
	}
	if oldStatus.IsAguardandoAprovacao() && request.ServiceOrderStatus.IsEmDiagnostico() {
		update.ServiceOrderStatus = valueobject.StatusEmDiagnostico
		// If the status is "EmDiagnostico", we can reset the PartsSupplies reserve
		for _, ps := range request.PartsSupplies {
			err := unreservePartsSupply(ctx, ps, partsSupplyRepo)
			if err != nil {
				log.Error().Msgf("Error unreserving parts supply: %v", err)
				return nil, err
			}
		}
		return update, nil
	}
	return nil, ErrInvalidTransitionStatusToEstimate
}

func ValidateExecution(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	oldStatus := serviceOrderDto.ServiceOrderStatus.ToDomain()

	if !request.ServiceOrderStatus.IsValid() {
		return nil, ErrInvalidStatus
	}
	if oldStatus.IsAprovada() && request.ServiceOrderStatus.IsEmExecucao() {
		update.ServiceOrderStatus = valueobject.StatusEmExecucao
		return update, nil
	}
	if oldStatus.IsEmExecucao() && request.ServiceOrderStatus.IsFinalizada() {
		update.ServiceOrderStatus = valueobject.StatusFinalizada
		return update, nil
	}
	return nil, ErrInvalidTransitionStatusToExecution
}

func ValidateDelivery(ctx context.Context, request *entities.ServiceOrder, serviceOrderDto *dto.ServiceOrderDTO, update *entities.ServiceOrder) (*entities.ServiceOrder, error) {
	oldStatus := serviceOrderDto.ServiceOrderStatus.ToDomain()

	if !request.ServiceOrderStatus.IsValid() {
		return nil, ErrInvalidStatus
	}

	if oldStatus.IsFinalizada() && request.ServiceOrderStatus.IsEntregue() {
		if serviceOrderDto.Payment == nil {
			log.Error().Msg("Payment information is required for delivery")
			return nil, errors.New("payment information is required for delivery")
		}
		update.ServiceOrderStatus = valueobject.StatusEntregue
		return update, nil
	}
	return nil, ErrInvalidTransitionStatusToDelivery
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

func getSeviceById(ctx context.Context, s entities.Service, serviceRepo service.IServiceRepo) (*entities.Service, error) {
	if s.ID == 0 {
		return nil, ErrInvalidID
	}
	result, err := serviceRepo.GetByID(ctx, s.ID)
	if err != nil {
		log.Error().Msgf("error finding service with id %d: %v", s.ID, err)
		return nil, err
	}
	if result.ID == 0 {
		log.Error().Msgf("service with id %d not found", s.ID)
		return nil, ErrServiceNotFound
	}
	return &result, nil
}

func getPartsSupplyByID(ctx context.Context, id uint, partsSupplyRepo parts_supply.IPartsSupplyRepo) (*entities.PartsSupply, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}
	result, err := partsSupplyRepo.GetByID(ctx, id)
	if err != nil {
		log.Error().Msgf("error finding parts supply with id %d: %v", id, err)
		return nil, err
	}
	if result.ID == 0 {
		log.Error().Msgf("parts supply with id %d not found", id)
		return nil, ErrPartsSupplyNotFound
	}
	return &result, nil
}

func validateQttPartsSupply(ctx context.Context, partsSupply entities.PartsSupply, partsSupplyRepo parts_supply.IPartsSupplyRepo) error {
	current, err := getPartsSupplyByID(ctx, partsSupply.ID, partsSupplyRepo)
	if err != nil {
		log.Error().Msgf("error getting parts supply by ID: %v", err)
		return err
	}

	totalAvailable := current.QuantityTotal - current.QuantityReserve
	if (partsSupply.QuantityReserve > totalAvailable) || (partsSupply.QuantityTotal > totalAvailable) {
		log.Error().Msgf("parts supply with id %d has insufficient quantity available", partsSupply.ID)
		return ErrInsufficientPartsSupply
	}
	return nil
}

func reservePartsSupply(ctx context.Context, partsSupply entities.PartsSupply, partsSupplyRepo parts_supply.IPartsSupplyRepo) error {
	current, err := getPartsSupplyByID(ctx, partsSupply.ID, partsSupplyRepo)
	if err != nil {
		log.Error().Msgf("error getting parts supply by ID: %v", err)
		return err
	}

	if partsSupply.QuantityReserve > 0 {
		current.QuantityReserve += partsSupply.QuantityReserve
	} else if partsSupply.QuantityTotal > 0 {
		current.QuantityTotal += partsSupply.QuantityTotal
	} else {
		return errors.New("no quantity to reserve")
	}

	err = partsSupplyRepo.Update(ctx, current)
	if err != nil {
		log.Error().Msgf("error reserving parts supply with id %d: %v", current.ID, err)
		return err
	}
	log.Info().Msgf("Parts supply with id %d reserved successfully", current.ID)
	return nil
}

// releaseReservedPartsSupply is when a service order is approved - Baixa de estoque
func releaseReservedPartsSupply(ctx context.Context, request entities.PartsSupply, partsSupplyRepo parts_supply.IPartsSupplyRepo) error {
	current, err := getPartsSupplyByID(ctx, request.ID, partsSupplyRepo)
	if err != nil {
		log.Error().Msgf("error getting parts supply by ID: %v", err)
		return err
	}

	if request.QuantityReserve > 0 {
		if current.QuantityReserve < request.QuantityReserve {
			return errors.New("cannot release more than reserved")
		}
		if request.QuantityReserve > 0 {
			current.QuantityReserve -= request.QuantityReserve
			current.QuantityTotal -= request.QuantityReserve
		} else if request.QuantityTotal > 0 {
			current.QuantityReserve -= request.QuantityTotal
			current.QuantityTotal -= request.QuantityTotal
		}
	} else {
		return errors.New("no quantity to release")
	}

	err = partsSupplyRepo.Update(ctx, current)
	if err != nil {
		log.Error().Msgf("error releasing reserved parts supply with id %d: %v", current.ID, err)
		return err
	}
	log.Info().Msgf("Reserved parts supply with id %d released successfully", current.ID)
	return nil
}

// unreservePartsSupply is when a service order is rejected - Liberação de reserva
func unreservePartsSupply(ctx context.Context, partsSupply entities.PartsSupply, partsSupplyRepo parts_supply.IPartsSupplyRepo) error {
	current, err := getPartsSupplyByID(ctx, partsSupply.ID, partsSupplyRepo)
	if err != nil {
		log.Error().Msgf("error getting parts supply by ID: %v", err)
		return err
	}

	if partsSupply.QuantityReserve > 0 {
		if current.QuantityReserve < partsSupply.QuantityReserve {
			return errors.New("cannot unreserve more than reserved")
		}
		if partsSupply.QuantityReserve > 0 {
			current.QuantityReserve -= partsSupply.QuantityReserve
		} else if partsSupply.QuantityTotal > 0 {
			current.QuantityReserve -= partsSupply.QuantityTotal
		}
	} else {
		return errors.New("no quantity to unreserved")
	}

	err = partsSupplyRepo.Update(ctx, current)
	if err != nil {
		log.Error().Msgf("error unreserving parts supply with id %d: %v", current.ID, err)
		return err
	}
	log.Info().Msgf("Parts supply with id %d unreserved successfully", current.ID)
	return nil
}

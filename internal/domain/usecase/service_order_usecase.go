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
	serviceorder "mecanica_xpto/internal/domain/repository/service_order"
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
	GetServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder) (*entities.ServiceOrder, error)
	ListServiceOrders(ctx context.Context) ([]*entities.ServiceOrder, error)
}

type ServiceOrderUseCase struct {
	repo            serviceorder.IServiceOrderRepository
	vehicleRepo     vehicles.VehicleRepositoryInterface
	customerRepo    customerRepo.ICustomerRepository
	serviceRepo     service.IServiceRepo
	partsSupplyRepo parts_supply.IPartsSupplyRepo
}

var _ IServiceOrderUseCase = (*ServiceOrderUseCase)(nil)

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

	update.ID = request.ID

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

	if (oldStatus.IsRecebida() || oldStatus.IsEmDiagnostico()) && newStatus.IsCancelada() {
		update.ServiceOrderStatus = valueobject.StatusCancelada
		return update, nil
	}

	if (oldStatus.IsRecebida() && newStatus.IsEmDiagnostico()) || oldStatus.IsEmDiagnostico() {
		update.ServiceOrderStatus = valueobject.StatusEmDiagnostico

		if len(request.Services) <= 0 && len(request.PartsSupplies) <= 0 {
			return update, nil
		}

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

		// TODO: a lista de services e peças está errada pois nao é as do banco de dados, é as do request
		// TODO: Depois validar como fazer um get dos services e parts supplies do banco de dados

		var err error

		// Calculate the total cost of PartsSupplies and Services and set it to the estimate
		update.Estimate, err = CalculateEstimate(ctx, update.Services, update.PartsSupplies, serviceRepo, partsSupplyRepo)
		if err != nil {
			log.Error().Msgf("Error calculating estimate: %v", err)
			return nil, err
		}

		log.Debug().Msgf("Estimate: %v", update.Estimate)

		// If OK, set a new status to "AguardandoAprovacao"
		update.ServiceOrderStatus = valueobject.StatusAguardandoAprovacao
		return update, nil
	}

	return nil, ErrInvalidTransitionStatusToDiagnosis

}

func CalculateEstimate(ctx context.Context, services []entities.Service, partsSupplies []entities.PartsSupply, serviceRepo service.IServiceRepo, psRepo parts_supply.IPartsSupplyRepo) (float64, error) {
	var totalEstimate float64

	servicesRegistered, err := getServicesByIDs(ctx, services, serviceRepo)
	if err != nil {
		log.Error().Msgf("Error getting services by IDs: %v", err)
		return 0, err
	}

	partsSuppliesRegistered, err := getPartsSupplyByIDs(ctx, partsSupplies, psRepo)
	if err != nil {
		log.Error().Msgf("Error getting parts supplies by IDs: %v", err)
		return 0, err
	}

	// Calculate services total
	for _, s := range servicesRegistered {
		log.Debug().Msgf("Service ID: %d, Price: %.2f", s.ID, s.Price)
		totalEstimate = totalEstimate + s.Price
	}

	log.Debug().Msgf("Total estimate from services: %.2f", totalEstimate)

	// Calculate parts supplies total
	for _, ps := range partsSuppliesRegistered {
		/* TODO: fazer a soma usando os preços de partsSuppliesRegistered e a quantidade de partsSupplies
		LOGS:
			{"level":"info","time":"2025-08-11T14:03:38-03:00","message":"Parts supply with id 1 reserved successfully"}
			{"level":"info","time":"2025-08-11T14:03:38-03:00","message":"Parts supply with id 2 reserved successfully"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"Service ID: 1, Price: 100.00"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"Service ID: 2, Price: 200.00"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"Total estimate from services: 300.00"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"PartsSupply ID: 1, Price: 50.00, QuantityReserve: 14, QuantityTotal: 100"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"PartsSupply ID: 2, Price: 100.00, QuantityReserve: 102, QuantityTotal: 200"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"Total estimate from parts supplies: 25300.00"}
			{"level":"debug","time":"2025-08-11T14:03:38-03:00","message":"Estimate: 25300"}
		*/
		log.Debug().Msgf("PartsSupply ID: %d, Price: %.2f, QuantityReserve: %d, QuantityTotal: %d", ps.ID, ps.Price, ps.QuantityReserve, ps.QuantityTotal)
		quantity := ps.QuantityReserve // Use reserve quantity by default
		if ps.QuantityTotal > 0 {      // If total quantity is specified, use that instead
			quantity = ps.QuantityTotal
		}
		totalEstimate = totalEstimate + (ps.Price * float64(quantity))
	}

	log.Debug().Msgf("Total estimate from parts supplies: %.2f", totalEstimate)

	return totalEstimate, nil
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
		current.QuantityReserve += partsSupply.QuantityTotal
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
		current.QuantityReserve -= request.QuantityReserve
		current.QuantityTotal -= request.QuantityReserve
	} else if request.QuantityTotal > 0 {
		if current.QuantityReserve < request.QuantityTotal {
			return errors.New("cannot release more than reserved")
		}
		current.QuantityReserve -= request.QuantityTotal
		current.QuantityTotal -= request.QuantityTotal
	} else {
		return errors.New("no quantity to release")
	}

	//TODO Remove logs when done
	log.Debug().Msgf("db qtt reserve: %d, qtt total: %d", current.QuantityReserve, current.QuantityTotal)
	log.Debug().Msgf("request qtt reserve: %d, qtt total: %d", request.QuantityReserve, request.QuantityTotal)

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
		return errors.New("no quantity to unreserve")
	}

	err = partsSupplyRepo.Update(ctx, current)
	if err != nil {
		log.Error().Msgf("error unreserving parts supply with id %d: %v", current.ID, err)
		return err
	}
	log.Info().Msgf("Parts supply with id %d unreserved successfully", current.ID)
	return nil
}

func (u *ServiceOrderUseCase) GetServiceOrder(ctx context.Context, serviceOrder entities.ServiceOrder) (*entities.ServiceOrder, error) {
	serviceOrderDto, err := u.repo.GetByID(serviceOrder.ID)
	if err != nil {
		log.Error().Msgf("error finding service order with id %d: %v", serviceOrder.ID, err)
		return nil, err
	}
	if serviceOrderDto == nil {
		log.Error().Msgf("service order with id %d not found", serviceOrder.ID)
		return nil, ErrServiceOrderNotFound
	}
	return serviceOrderDto.ToDomain(), nil
}

func (u *ServiceOrderUseCase) ListServiceOrders(ctx context.Context) ([]*entities.ServiceOrder, error) {
	var serviceOrdersResponse []*entities.ServiceOrder
	serviceOrders, err := u.repo.List()
	if err != nil {
		log.Error().Msgf("error listing service orders: %v", err)
		return nil, err
	}
	for _, so := range serviceOrders {
		serviceOrdersResponse = append(serviceOrdersResponse, so.ToDomain())
	}

	return serviceOrdersResponse, nil
}

func getServicesByIDs(ctx context.Context, services []entities.Service, serviceRepo service.IServiceRepo) ([]entities.Service, error) {
	if len(services) == 0 {
		return nil, errors.New("no services provided")
	}

	var serviceDb []entities.Service
	for _, s := range services {
		item, err := serviceRepo.GetByID(ctx, s.ID)
		if err != nil {
			log.Error().Msgf("error getting service by ID %v: %v", s.ID, err)
			return nil, err
		}
		serviceDb = append(serviceDb, item)
	}

	// Assuming we have a service repository to get the services by IDs
	return serviceDb, nil
}

func getPartsSupplyByIDs(ctx context.Context, partsSupplies []entities.PartsSupply, psRepo parts_supply.IPartsSupplyRepo) ([]entities.PartsSupply, error) {
	if len(partsSupplies) == 0 {
		return nil, errors.New("no services provided")
	}

	var psDb []entities.PartsSupply
	for _, ps := range partsSupplies {
		item, err := psRepo.GetByID(ctx, ps.ID)
		if err != nil {
			log.Error().Msgf("error getting Parts Supplies by ID %v: %v", ps.ID, err)
			return nil, err
		}
		psDb = append(psDb, item)
	}

	return psDb, nil
}

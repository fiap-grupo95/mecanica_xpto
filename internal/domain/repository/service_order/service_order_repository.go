package serviceorder

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"

	"gorm.io/gorm"
)

type IServiceOrderRepository interface {
	Create(serviceOrder *entities.ServiceOrder) error
	GetByID(id uint) (*dto.ServiceOrderDTO, error)
	Update(serviceOrder *entities.ServiceOrder) error
	List() ([]dto.ServiceOrderDTO, error)
	GetStatus(status valueobject.ServiceOrderStatus) (*dto.ServiceOrderStatusDTO, error)
}

// ServiceOrderRepository implements IServiceOrderRepository interface
type ServiceOrderRepository struct {
	db *gorm.DB
}

func NewServiceOrderRepository(db *gorm.DB) *ServiceOrderRepository {
	return &ServiceOrderRepository{db: db}
}

func (r *ServiceOrderRepository) Create(serviceOrder *entities.ServiceOrder) error {
	if serviceOrder == nil {
		return gorm.ErrInvalidData
	}

	dtoStatus, err := r.GetStatus(serviceOrder.ServiceOrderStatus)
	if err != nil {
		return gorm.ErrInvalidData
	}

	if dtoStatus == nil {
		return gorm.ErrInvalidData
	}

	// Begin transaction
	tx := r.db.Begin()

	serviceOrderDto := dto.ServiceOrderDTO{
		ID:                   serviceOrder.ID,
		CustomerID:           serviceOrder.CustomerID,
		VehicleID:            serviceOrder.VehicleID,
		OSStatusID:           dtoStatus.ID,
		Estimate:             serviceOrder.Estimate,
		StartedExecutionDate: serviceOrder.StartedExecutionDate,
		FinalExecutionDate:   serviceOrder.FinalExecutionDate,
		CreatedAt:            serviceOrder.CreatedAt,
		UpdatedAt:            serviceOrder.UpdatedAt,
	}

	if err := tx.Save(&serviceOrderDto).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Handle PartsSupplies N:N relationship
	for _, partsSupply := range serviceOrder.PartsSupplies {
		relation := dto.PartsSupplyServiceOrderDTO{
			PartsSupplyID:  partsSupply.ID,
			ServiceOrderID: serviceOrderDto.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Handle Services N:N relationship
	for _, service := range serviceOrder.Services {
		relation := dto.ServiceServiceOrderDTO{
			ServiceID:      service.ID,
			ServiceOrderID: serviceOrderDto.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *ServiceOrderRepository) GetByID(id uint) (*dto.ServiceOrderDTO, error) {
	var serviceOrder dto.ServiceOrderDTO
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessário
	err := r.db.Preload("Customer").
		Preload("Vehicle").
		Preload("ServiceOrderStatus").
		Preload("AdditionalRepairs").
		Preload("Payment").
		Preload("PartsSupplies").
		Preload("Services").
		First(&serviceOrder, id).Error
	if err != nil {
		return nil, err
	}
	return &serviceOrder, nil
}

func (r *ServiceOrderRepository) Update(serviceOrder *entities.ServiceOrder) error {
	if serviceOrder == nil {
		return gorm.ErrInvalidData
	}

	dtoStatus, err := r.GetStatus(serviceOrder.ServiceOrderStatus)
	if err != nil {
		return gorm.ErrInvalidData
	}

	tx := r.db.Begin()

	serviceOrderDto := dto.ServiceOrderDTO{
		ID:                   serviceOrder.ID,
		CustomerID:           serviceOrder.CustomerID,
		VehicleID:            serviceOrder.VehicleID,
		OSStatusID:           dtoStatus.ID,
		Estimate:             serviceOrder.Estimate,
		StartedExecutionDate: serviceOrder.StartedExecutionDate,
		FinalExecutionDate:   serviceOrder.FinalExecutionDate,
	}

	if err := tx.Model(&dto.ServiceOrderDTO{}).Where("id = ?", serviceOrder.ID).Updates(&serviceOrderDto).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update PartsSupplies relationships
	if err := tx.Where("service_order_id = ?", serviceOrder.ID).Delete(&dto.PartsSupplyServiceOrderDTO{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, partsSupply := range serviceOrder.PartsSupplies {
		relation := dto.PartsSupplyServiceOrderDTO{
			PartsSupplyID:  partsSupply.ID,
			ServiceOrderID: serviceOrder.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Services relationships
	if err := tx.Where("service_order_id = ?", serviceOrder.ID).Delete(&dto.ServiceServiceOrderDTO{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, service := range serviceOrder.Services {
		relation := dto.ServiceServiceOrderDTO{
			ServiceID:      service.ID,
			ServiceOrderID: serviceOrder.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *ServiceOrderRepository) List() ([]dto.ServiceOrderDTO, error) {
	var serviceOrders []dto.ServiceOrderDTO
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessário
	err := r.db.Preload("Customer").
		Preload("Vehicle").
		Preload("ServiceOrderStatus").
		Preload("AdditionalRepairs").
		Preload("Payment").
		Preload("PartsSupplies").
		Preload("Services").
		Find(&serviceOrders).Error
	return serviceOrders, err
}

func (r *ServiceOrderRepository) GetStatus(status valueobject.ServiceOrderStatus) (*dto.ServiceOrderStatusDTO, error) {
	var serviceOrderStatuses dto.ServiceOrderStatusDTO
	err := r.db.Where("description = ?", status.String()).First(&serviceOrderStatuses).Error
	if err != nil {
		return nil, err
	}
	return &serviceOrderStatuses, nil
}

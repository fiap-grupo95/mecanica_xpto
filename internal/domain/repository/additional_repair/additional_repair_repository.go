package additional_repair

import (
	"gorm.io/gorm"
	"mecanica_xpto/internal/domain/model/dto"
)

type IAdditionalRepairRepository interface {
	Create(additionalRepair *dto.AdditionalRepairDTO) error
	GetByID(id uint) (*dto.AdditionalRepairDTO, error)
	AddPartSupplyAndService(id uint, updatedAdditionalRepair *dto.AdditionalRepairDTO) error
	RemovePartSupplyAndService(id uint, updatedAdditionalRepair *dto.AdditionalRepairDTO) error
	GetByServiceOrder(serviceOrderId uint) ([]dto.AdditionalRepairDTO, error)
	GetStatus(status string) (*dto.AdditionalRepairStatusDTO, error)
}

// AdditionalRepairRepository implements IAdditionalRepairRepository interface
type AdditionalRepairRepository struct {
	db *gorm.DB
}

func NewAdditionalRepairRepository(db *gorm.DB) IAdditionalRepairRepository {
	return &AdditionalRepairRepository{db: db}
}

func (r *AdditionalRepairRepository) Create(additionalRepair *dto.AdditionalRepairDTO) error {
	dtoStatus, err := r.GetStatus(additionalRepair.ARStatus.Description)
	if err != nil {
		return gorm.ErrInvalidData
	}
	additionalRepair.ARStatus = *dtoStatus
	tx := r.db.Begin()
	if err := tx.Create(&additionalRepair).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) GetByID(id uint) (*dto.AdditionalRepairDTO, error) {
	var additionalRepair dto.AdditionalRepairDTO
	err := r.db.Preload("ARStatus").
		Preload("PartsSupplies").
		Preload("Services").
		First(&additionalRepair, id).Error
	if err != nil {
		return nil, err
	}
	return &additionalRepair, nil
}

func (r *AdditionalRepairRepository) AddPartSupplyAndService(id uint, updatedAdditionalRepair *dto.AdditionalRepairDTO) error {
	var dtoDB dto.AdditionalRepairDTO
	if err := r.db.First(&dtoDB, id).Error; err != nil {
		return err
	}
	tx := r.db.Begin()

	for _, ps := range updatedAdditionalRepair.PartsSupplies {
		relation := dto.PartsSupplyAdditionalRepairDTO{
			PartsSupplyID:      ps.ID,
			AdditionalRepairID: id,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, svc := range updatedAdditionalRepair.Services {
		relation := dto.ServiceAdditionalRepairDTO{
			ServiceID:          svc.ID,
			AdditionalRepairID: id,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	newEstimate := calculateEstimate(updatedAdditionalRepair.Services, updatedAdditionalRepair.PartsSupplies)

	// Update only the Estimate field
	if err := tx.Model(&dto.AdditionalRepairDTO{}).
		Where("id = ?", id).
		Update("estimate", newEstimate).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) RemovePartSupplyAndService(id uint, updatedAdditionalRepair *dto.AdditionalRepairDTO) error {
	var dtoDB dto.AdditionalRepairDTO
	if err := r.db.First(&dtoDB, id).Error; err != nil {
		return err
	}
	tx := r.db.Begin()

	if err := tx.Where("additional_repair_id = ?", id).
		Delete(&dto.PartsSupplyAdditionalRepairDTO{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("additional_repair_id = ?", id).
		Delete(&dto.ServiceAdditionalRepairDTO{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	newEstimate := recalculateEstimateAfterRemoval(&dtoDB, updatedAdditionalRepair.PartsSupplies, updatedAdditionalRepair.Services)

	// Update only the Estimate field
	if err := tx.Model(&dto.AdditionalRepairDTO{}).
		Where("id = ?", id).
		Update("estimate", newEstimate).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) GetByServiceOrder(serviceOrderId uint) ([]dto.AdditionalRepairDTO, error) {
	var additionalRepairs []dto.AdditionalRepairDTO
	err := r.db.Preload("ARStatus").
		Preload("PartsSupplies").
		Preload("Services").
		Where("service_order_id = ?", serviceOrderId).
		Find(&additionalRepairs).Error
	if err != nil {
		return nil, err
	}
	return additionalRepairs, nil
}

func (r *AdditionalRepairRepository) GetStatus(status string) (*dto.AdditionalRepairStatusDTO, error) {
	var additionalRepairStatus dto.AdditionalRepairStatusDTO
	err := r.db.Where("description = ?", status).First(&additionalRepairStatus).Error
	if err != nil {
		return nil, err
	}
	return &additionalRepairStatus, nil
}

func calculateEstimate(services []dto.ServiceDTO, partsSupplies []dto.PartsSupplyDTO) float64 {
	var total float64
	for _, svc := range services {
		total += svc.Price
	}
	for _, ps := range partsSupplies {
		total += ps.Price
	}
	return total
}

func recalculateEstimateAfterRemoval(additionalRepair *dto.AdditionalRepairDTO, removedPartsSupplies []dto.PartsSupplyDTO, removedServices []dto.ServiceDTO) float64 {
	estimate := additionalRepair.Estimate
	for _, svc := range removedServices {
		estimate -= svc.Price
	}
	for _, ps := range removedPartsSupplies {
		estimate -= ps.Price
	}
	if estimate < 0 {
		estimate = 0
	}
	return estimate
}

package additional_repair

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"time"

	"gorm.io/gorm"
)

type IAdditionalRepairRepository interface {
	Create(additionalRepair *entities.AdditionalRepair) error
	GetByID(id uint) (*dto.AdditionalRepairDTO, error)
	Update(additionalRepair *entities.AdditionalRepair) error
	GetByServiceOrder(serviceOrderId int) ([]dto.AdditionalRepairDTO, error)
	GetStatus(status string) (*dto.AdditionalRepairStatusDTO, error)
}

// AdditionalRepairRepository implements IAdditionalRepairRepository interface
type AdditionalRepairRepository struct {
	db *gorm.DB
}

func NewAdditionalRepairRepository(db *gorm.DB) IAdditionalRepairRepository {
	return &AdditionalRepairRepository{db: db}
}

func (r *AdditionalRepairRepository) Create(additionalRepair *entities.AdditionalRepair) error {
	arStatus, err := r.GetStatus(string(additionalRepair.ARStatus))
	if err != nil {
		return gorm.ErrInvalidData
	}

	if arStatus == nil {
		return gorm.ErrInvalidData
	}

	additionalRepairDto := dto.AdditionalRepairDTO{
		ID:         additionalRepair.ID,
		ARStatusID: arStatus.ID,
		Estimate:   additionalRepair.Estimate,
		CreatedAt:  additionalRepair.CreatedAt,
		ID:            additionalRepair.ID,
		ARStatusID:    arStatus.ID,
		Estimate:      additionalRepair.Estimate,
		ServiceOrderID: additionalRepair.ServiceOrderID,
		CreatedAt:     additionalRepair.CreatedAt,
		UpdatedAt:     additionalRepair.UpdatedAt,
	}

	tx := r.db.Begin()
	if err := tx.Create(&additionalRepairDto).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Handle PartsSupplies N:N relationship
	for _, partsSupply := range additionalRepair.PartsSupplies {
		relation := dto.PartsSupplyAdditionalRepairDTO{
			PartsSupplyID:      partsSupply.ID,
			AdditionalRepairID: additionalRepairDto.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Handle Services N:N relationship
	for _, service := range additionalRepair.Services {
		relation := dto.ServiceAdditionalRepairDTO{
			ServiceID:          service.ID,
			AdditionalRepairID: additionalRepairDto.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) GetByID(id uint) (*dto.AdditionalRepairDTO, error) {
	var additionalRepair dto.AdditionalRepairDTO
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessário
	err := r.db.Preload("ARStatus").
		Preload("PartsSupplies").
		Preload("Services").
		First(&additionalRepair, id).Error
	if err != nil {
		return nil, err
	}
	return &additionalRepair, nil
}

func (r *AdditionalRepairRepository) Update(additionalRepair *entities.AdditionalRepair) error {
	arStatus, err := r.GetStatus(string(additionalRepair.ARStatus))
	if err != nil {
		return gorm.ErrInvalidData
	}

	tx := r.db.Begin()

	// Update main AdditionalRepair record
	additionalRepairDto := dto.AdditionalRepairDTO{
		ID:         additionalRepair.ID,
		ARStatusID: arStatus.ID,
		Estimate:   additionalRepair.Estimate,
		UpdatedAt:  time.Now(),
	}

	if err := tx.Model(&dto.AdditionalRepairDTO{}).Where("id = ?", additionalRepair.ID).Updates(&additionalRepairDto).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing relationships
	if err := tx.Where("additional_repair_id = ?", additionalRepair.ID).Delete(&dto.PartsSupplyAdditionalRepairDTO{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("additional_repair_id = ?", additionalRepair.ID).Delete(&dto.ServiceAdditionalRepairDTO{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recreate PartsSupplies relationships
	for _, partsSupply := range additionalRepair.PartsSupplies {
		relation := dto.PartsSupplyAdditionalRepairDTO{
			PartsSupplyID:      partsSupply.ID,
			AdditionalRepairID: additionalRepair.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Recreate Services relationships
	for _, service := range additionalRepair.Services {
		relation := dto.ServiceAdditionalRepairDTO{
			ServiceID:          service.ID,
			AdditionalRepairID: additionalRepair.ID,
		}
		if err := tx.Create(&relation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *AdditionalRepairRepository) GetByServiceOrder(serviceOrderId int) ([]dto.AdditionalRepairDTO, error) {
	var additionalRepairs []dto.AdditionalRepairDTO
	// TODO - Avaliar o que posso tirar do Preload e deixar para serem carregados apenas quando necessário
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

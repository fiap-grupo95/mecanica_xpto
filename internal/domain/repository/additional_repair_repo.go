package repository

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
)

// 1:N relationship between AdditionalRepair and AdditionalRepairStatus
type AdditionalRepairStatusDTO struct {
	ID                uint   `gorm:"primaryKey"`
	Description       string `gorm:"size:50;not null"`
	AdditionalRepairs []AdditionalRepairDTO
}

func (arsm AdditionalRepairStatusDTO) ToDomain() entities.AdditionalRepairStatus {
	return entities.AdditionalRepairStatus{
		ID:                arsm.ID,
		Description:       arsm.Description,
		AdditionalRepairs: nil, // This will be populated later if needed
	}
}

type AdditionalRepairDTO struct {
	ID             uint                      `gorm:"primaryKey"`
	ServiceOrderID uint                      `gorm:"not null"`
	ServiceOrder   ServiceOrderDTO           `gorm:"foreignKey:ServiceOrderID"`
	ServiceID      uint                      `gorm:"not null"`
	Service        dto.ServiceDTO            `gorm:"foreignKey:ServiceID"`
	PartsSupplyID  uint                      `gorm:"not null"`
	PartsSupply    dto.PartsSupplyDTO        `gorm:"foreignKey:PartsSupplyID"`
	ARStatusID     uint                      `gorm:"not null"`
	ARStatus       AdditionalRepairStatusDTO `gorm:"foreignKey:ARStatusID"`
}

func (arm *AdditionalRepairDTO) ToDomain() entities.AdditionalRepair {
	return entities.AdditionalRepair{
		ID:           arm.ID,
		ServiceOrder: arm.ServiceOrder.ToDomain(),
		Service:      arm.Service.ToDomain(),
		PartsSupply:  arm.PartsSupply.ToDomain(),
		ARStatus:     valueobject.ParseAdditionalRepairStatus(arm.ARStatus.Description),
	}
}

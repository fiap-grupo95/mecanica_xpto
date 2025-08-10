package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

// N:N relationship between PartsSupply and AdditionalRepair
type PartsSupplyAdditionalRepairDTO struct {
	PartsSupplyID      uint `gorm:"primaryKey"`
	AdditionalRepairID uint `gorm:"primaryKey"`
}

// N:N relationship between Service and AdditionalRepair
type ServiceAdditionalRepairDTO struct {
	ServiceID          uint `gorm:"primaryKey"`
	AdditionalRepairID uint `gorm:"primaryKey"`
}

// 1:N relationship between AdditionalRepair and AdditionalRepairStatus
type AdditionalRepairStatusDTO struct {
	ID                uint                  `gorm:"primaryKey"`
	Description       string                `gorm:"size:50;not null"`
	AdditionalRepairs []AdditionalRepairDTO `gorm:"foreignKey:ARStatusID"`
}

func (arsm AdditionalRepairStatusDTO) ToDomain() valueobject.AdditionalRepairStatus {
	return valueobject.ParseAdditionalRepairStatus(arsm.Description)
}

type AdditionalRepairDTO struct {
	ID             uint                      `gorm:"primaryKey"`
	ServiceOrderID uint                      `gorm:"not null"`
	ServiceOrder   ServiceOrderDTO           `gorm:"foreignKey:ServiceOrderID"`
	ServiceID      uint                      `gorm:"column:service_order_id;not null"`
	ARStatusID     uint                      `gorm:"not null"`
	ARStatus       AdditionalRepairStatusDTO `gorm:"foreignKey:ARStatusID"`
	Estimate       float64                   `gorm:"type:decimal(10,2)"`
	CreatedAt      time.Time                 `gorm:"autoCreateTime"`
	UpdatedAt      time.Time                 `gorm:"autoUpdateTime"`
	PartsSupplies  []PartsSupplyDTO          `gorm:"foreignKey:AdditionalRepairID"`
	Services       []ServiceDTO              `gorm:"foreignKey:AdditionalRepairID"`
}

func (arm *AdditionalRepairDTO) ToDomain() entities.AdditionalRepair {
	var partsSupplies []entities.PartsSupply
	var services []entities.Service

	for _, ps := range arm.PartsSupplies {
		partsSupplies = append(partsSupplies, ps.ToDomain())
	}

	for _, s := range arm.Services {
		services = append(services, s.ToDomain())
	}

	return entities.AdditionalRepair{
		ID:            arm.ID,
		ARStatus:      arm.ARStatus.ToDomain(),
		Estimate:      arm.Estimate,
		CreatedAt:     arm.CreatedAt,
		UpdatedAt:     arm.UpdatedAt,
		PartsSupplies: partsSupplies,
		Services:      services,
	}
}

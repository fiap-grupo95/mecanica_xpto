package repository

import "mecanica_xpto/internal/domain/entities"

// 1:N relationship between AdditionalRepair and AdditionalRepairStatus
type AdditionalRepairStatusModel struct {
	ID                uint   `gorm:"primaryKey"`
	Description       string `gorm:"size:50;not null"`
	AdditionalRepairs []AdditionalRepairModel
}

// crie um toDomain para a entidade AdditionalRepairStatusModel
func (arsm AdditionalRepairStatusModel) ToDomain() entities.AdditionalRepairStatus {
	return entities.AdditionalRepairStatus{
		ID:                arsm.ID,
		Description:       arsm.Description,
		AdditionalRepairs: nil, // This will be populated later if needed
	}
}

type AdditionalRepairModel struct {
	ID             uint                        `gorm:"primaryKey"`
	ServiceOrderID uint                        `gorm:"not null"`
	ServiceOrder   ServiceOrderModel           `gorm:"foreignKey:ServiceOrderID"`
	ServiceID      uint                        `gorm:"not null"`
	Service        ServiceModel                `gorm:"foreignKey:ServiceID"`
	PartsSupplyID  uint                        `gorm:"not null"`
	PartsSupply    PartsSupplyModel            `gorm:"foreignKey:PartsSupplyID"`
	ARStatusID     uint                        `gorm:"not null"`
	ARStatus       AdditionalRepairStatusModel `gorm:"foreignKey:ARStatusID"`
}

// crie um toDomain para a entidade AdditionalRepairModel
func (arm *AdditionalRepairModel) ToDomain() entities.AdditionalRepair {
	return entities.AdditionalRepair{
		ID:             arm.ID,
		ServiceOrderID: arm.ServiceOrderID,
		ServiceOrder:   arm.ServiceOrder.ToDomain(),
		ServiceID:      arm.ServiceID,
		Service:        arm.Service.ToDomain(),
		PartsSupplyID:  arm.PartsSupplyID,
		PartsSupply:    arm.PartsSupply.ToDomain(),
		ARStatusID:     arm.ARStatusID,
		ARStatus:       arm.ARStatus.ToDomain(),
	}
}

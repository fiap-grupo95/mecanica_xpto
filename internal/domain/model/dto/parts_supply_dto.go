package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"time"

	"gorm.io/gorm"
)

// N:N relationship between PartsSupply and ServiceOrder
// 1:N relationship between PartsSupply and AdditionalRepair
type PartsSupplyDTO struct {
	ID                uint                  `gorm:"primaryKey"`
	Name              string                `gorm:"size:100;not null"`
	Description       string                `gorm:"type:text"`
	Price             float64               `gorm:"type:decimal(10,2);not null"`
	QuantityTotal     int                   `gorm:"not null;default:0"`
	QuantityReserve   int                   `gorm:"not null;default:0"`
	CreatedAt         time.Time             `gorm:"autoCreateTime"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt        `gorm:"index"`
	AdditionalRepairs []AdditionalRepairDTO `gorm:"many2many:parts_supply_additional_repair"`
	ServiceOrders     []ServiceOrderDTO     `gorm:"many2many:parts_supply_service_order_dtos;joinForeignKey:parts_supply_id;joinReferences:service_order_id"`
}

func (m *PartsSupplyDTO) ToDomain() entities.PartsSupply {
	return entities.PartsSupply{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		Price:           m.Price,
		QuantityTotal:   m.QuantityTotal,
		QuantityReserve: m.QuantityReserve,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt: func() *time.Time {
			if m.DeletedAt.Valid {
				return &m.DeletedAt.Time
			}
			return nil
		}(),
		AdditionalRepairs: nil, // This will be populated by the repository layer
		ServiceOrders:     nil, // This will be populated by the repository layer
	}
}

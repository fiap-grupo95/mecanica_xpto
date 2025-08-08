package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"time"

	"gorm.io/gorm"
)

type ServiceDTO struct {
	ID                uint                  `gorm:"primaryKey"`
	Name              string                `gorm:"size:100;not null"`
	Description       string                `gorm:"type:text"`
	Price             float64               `gorm:"type:decimal(10,2);not null"`
	CreatedAt         time.Time             `gorm:"autoCreateTime"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt        `gorm:"index"`
	AdditionalRepairs []AdditionalRepairDTO `gorm:"foreignKey:ServiceID"`
	ServiceOrders     []ServiceOrderDTO     `gorm:"many2many:tb_service_service_orders;"`
}

func (m *ServiceDTO) ToDomain() entities.Service {
	return entities.Service{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt: func() *time.Time {
			if m.DeletedAt.Valid {
				return &m.DeletedAt.Time
			}
			return nil
		}(),
	}
}

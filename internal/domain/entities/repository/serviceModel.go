package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type ServiceModel struct {
	ID                uint           `gorm:"primaryKey"`
	Name              string         `gorm:"size:100;not null"`
	Description       string         `gorm:"type:text"`
	Price             float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	AdditionalRepairs []AdditionalRepairModel
	ServiceOrders     []ServiceOrderModel `gorm:"many2many:service_serviceorders;"`
}

func (m *ServiceModel) ToDomain() entities.Service {
	return entities.Service{
		ID:                m.ID,
		Name:              m.Name,
		Description:       m.Description,
		Price:             m.Price,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		DeletedAt:         nil, // This will be populated by the repository layer
		AdditionalRepairs: nil, // This will be populated by the repository layer
		ServiceOrders:     nil, // This will be populated by the repository layer
	}
}

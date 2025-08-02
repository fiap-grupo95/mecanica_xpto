package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/entities/valueobject"
	"time"
)

type ServiceOrderStatusModel struct {
	ID            uint   `gorm:"primaryKey"`
	Description   string `gorm:"size:50;not null"`
	ServiceOrders []ServiceOrderModel
}

type ServiceOrderModel struct {
	ID                   uint                    `gorm:"primaryKey"`
	CustomerID           uint                    `gorm:"not null"`
	Customer             CustomerModel           `gorm:"foreignKey:CustomerID"`
	VehicleID            uint                    `gorm:"not null"`
	Vehicle              VehicleModel            `gorm:"foreignKey:VehicleID"`
	OSStatusID           uint                    `gorm:"not null"`
	ServiceOrderStatus   ServiceOrderStatusModel `gorm:"foreignKey:OSStatusID"`
	Estimate             float64                 `gorm:"type:decimal(10,2)"`
	StartedExecutionDate time.Time
	FinalExecutionDate   time.Time
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	AdditionalRepairs    []AdditionalRepairModel
	Payment              *PaymentModel      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PartsSupplies        []PartsSupplyModel `gorm:"many2many:partssupply_serviceorders;"`
	Services             []ServiceModel     `gorm:"many2many:service_serviceorders;"`
}

func (m *ServiceOrderStatusModel) ToDomain() entities.ServiceOrderStatus {
	return entities.ServiceOrderStatus{
		ID:          m.ID,
		Description: valueobject.ParseServiceOrderStatus(m.Description),
	}
}

func (m *ServiceOrderModel) ToDomain() entities.ServiceOrder {
	return entities.ServiceOrder{
		ID:                   m.ID,
		Estimate:             m.Estimate,
		StartedExecutionDate: m.StartedExecutionDate,
		FinalExecutionDate:   m.FinalExecutionDate,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
		ServiceOrderStatus:   valueobject.ParseServiceOrderStatus(m.ServiceOrderStatus.Description),
		AdditionalRepairs:    nil, // This will be populated by the repository layer
		Payment:              nil, // This will be populated by the repository layer
		PartsSupplies:        nil, // This will be populated by the repository layer
		Services:             nil, // This will be populated by the repository layer
	}
}

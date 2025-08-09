package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

type ServiceOrderStatusDTO struct {
	ID            uint              `gorm:"primaryKey"`
	Description   string            `gorm:"size:50;not null"`
	ServiceOrders []ServiceOrderDTO `gorm:"foreignKey:OSStatusID"`
}

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrderDTO struct {
	PartsSupplyID  uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
}

// N:N relationship between Service and ServiceOrder
type ServiceServiceOrderDTO struct {
	ServiceID      uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
}

type ServiceOrderDTO struct {
	ID                   uint                  `gorm:"primaryKey"`
	CustomerID           uint                  `gorm:"not null"`
	Customer             CustomerDTO           `gorm:"foreignKey:CustomerID"`
	VehicleID            uint                  `gorm:"not null"`
	Vehicle              VehicleDTO            `gorm:"foreignKey:VehicleID"`
	OSStatusID           uint                  `gorm:"not null"`
	ServiceOrderStatus   ServiceOrderStatusDTO `gorm:"foreignKey:OSStatusID"`
	Estimate             float64               `gorm:"type:decimal(10,2)"`
	StartedExecutionDate time.Time
	FinalExecutionDate   time.Time
	CreatedAt            time.Time             `gorm:"autoCreateTime"`
	UpdatedAt            time.Time             `gorm:"autoUpdateTime"`
	AdditionalRepairs    []AdditionalRepairDTO `gorm:"foreignKey:ServiceOrderID"`
	Payment              *PaymentDTO           `gorm:"foreignKey:ServiceOrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PartsSupplies        []PartsSupplyDTO      `gorm:"many2many:parts_supply_service_order_dtos;"`
	Services             []ServiceDTO          `gorm:"many2many:service_service_order_dtos;"`
}

func (m *ServiceOrderStatusDTO) ToDomain() entities.ServiceOrderStatus {
	return entities.ServiceOrderStatus{
		ID:          m.ID,
		Description: valueobject.ParseServiceOrderStatus(m.Description),
	}
}

func (m *ServiceOrderDTO) ToDomain() entities.ServiceOrder {
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

package repository

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

type ServiceOrderStatusDTO struct {
	ID            uint   `gorm:"primaryKey"`
	Description   string `gorm:"size:50;not null"`
	ServiceOrders []ServiceOrderDTO
}

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrderDTO struct {
	PartsSupplyID  uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
	Quantity       int  `gorm:"not null;default:1"`
}

// N:N relationship between Service and ServiceOrder
type ServiceServiceOrderDTO struct {
	ServiceID      uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
}

type ServiceOrderDTO struct {
	ID                   uint                  `gorm:"primaryKey"`
	CustomerID           uint                  `gorm:"not null"`
	Customer             dto.CustomerDTO       `gorm:"foreignKey:CustomerID"`
	VehicleID            uint                  `gorm:"not null"`
	Vehicle              dto.VehicleDTO        `gorm:"foreignKey:VehicleID"`
	OSStatusID           uint                  `gorm:"not null"`
	ServiceOrderStatus   ServiceOrderStatusDTO `gorm:"foreignKey:OSStatusID"`
	Estimate             float64               `gorm:"type:decimal(10,2)"`
	StartedExecutionDate time.Time
	FinalExecutionDate   time.Time
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	AdditionalRepairs    []AdditionalRepairDTO
	Payment              *PaymentDTO          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PartsSupplies        []dto.PartsSupplyDTO `gorm:"many2many:partssupply_serviceorders;"`
	Services             []dto.ServiceDTO     `gorm:"many2many:service_serviceorders;"`
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

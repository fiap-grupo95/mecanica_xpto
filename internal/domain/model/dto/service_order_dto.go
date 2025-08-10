package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"
)

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

type ServiceOrderStatusDTO struct {
	ID          uint   `gorm:"primaryKey"`
	Description string `gorm:"size:50;not null"`
}

func (m *ServiceOrderStatusDTO) ToDomain() valueobject.ServiceOrderStatus {
	return valueobject.ParseServiceOrderStatus(m.Description)
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
	PartsSupplies        []PartsSupplyDTO      `gorm:"many2many:parts_supply_service_order;"`
	Services             []ServiceDTO          `gorm:"many2many:service_service_order;"`
}

func (m *ServiceOrderDTO) ToDomain() *entities.ServiceOrder {
	var additionalRepairs []entities.AdditionalRepair
	var partsSupplies []entities.PartsSupply
	var services []entities.Service

	// Convert AdditionalRepairs
	for _, ar := range m.AdditionalRepairs {
		additionalRepairs = append(additionalRepairs, ar.ToDomain())
	}

	// Convert PartsSupplies
	for _, ps := range m.PartsSupplies {
		partsSupplies = append(partsSupplies, ps.ToDomain())
	}

	// Convert Services
	for _, s := range m.Services {
		services = append(services, s.ToDomain())
	}

	// Convert Payment if exists
	var payment *entities.Payment
	if m.Payment != nil {
		p := m.Payment.ToDomain()
		payment = &p
	}

	// Convert Customer and Vehicle if they are loaded
	var customer entities.Customer
	var vehicle entities.Vehicle
	if m.Customer.ID != 0 {
		customer = *m.Customer.ToDomain()
	}
	if m.Vehicle.ID != 0 {
		vehicle = *m.Vehicle.ToDomain()
	}

	return &entities.ServiceOrder{
		ID:                   m.ID,
		CustomerID:           m.CustomerID,
		Customer:             customer,
		VehicleID:            m.VehicleID,
		Vehicle:              vehicle,
		ServiceOrderStatus:   m.ServiceOrderStatus.ToDomain(),
		Estimate:             m.Estimate,
		StartedExecutionDate: m.StartedExecutionDate,
		FinalExecutionDate:   m.FinalExecutionDate,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
		AdditionalRepairs:    additionalRepairs,
		Payment:              payment,
		PartsSupplies:        partsSupplies,
		Services:             services,
	}
}

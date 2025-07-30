package entities

import "time"

type ServiceOrderStatus struct {
	ID            uint   `gorm:"primaryKey"`
	Description   string `gorm:"size:50;not null"`
	ServiceOrders []ServiceOrder
}

type ServiceOrder struct {
	ID                   uint               `gorm:"primaryKey"`
	CustomerID           uint               `gorm:"not null"`
	Customer             Customer           `gorm:"foreignKey:CustomerID"`
	VehicleID            uint               `gorm:"not null"`
	Vehicle              Vehicle            `gorm:"foreignKey:VehicleID"`
	OSStatusID           uint               `gorm:"not null"`
	ServiceOrderStatus   ServiceOrderStatus `gorm:"foreignKey:OSStatusID"`
	Estimate             float64            `gorm:"type:decimal(10,2)"`
	StartedExecutionDate time.Time
	FinalExecutionDate   time.Time
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	AdditionalRepairs    []AdditionalRepair
	Payment              *Payment      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PartsSupplies        []PartsSupply `gorm:"many2many:partssupply_serviceorders;"`
	Services             []Service     `gorm:"many2many:service_serviceorders;"`
}

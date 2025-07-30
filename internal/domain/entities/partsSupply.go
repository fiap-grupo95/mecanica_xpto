package entities

import (
	"time"

	"gorm.io/gorm"
)

// N:N relationship between PartsSupply and ServiceOrder
// 1:N relationship between PartsSupply and AdditionalRepair
type PartsSupply struct {
	ID                uint           `gorm:"primaryKey"`
	Name              string         `gorm:"size:100;not null"`
	Description       string         `gorm:"type:text"`
	Price             float64        `gorm:"type:decimal(10,2);not null"`
	QuantityTotal     int            `gorm:"not null;default:0"`
	QuantityReserve   int            `gorm:"not null;default:0"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	AdditionalRepairs []AdditionalRepair
	ServiceOrders     []ServiceOrder `gorm:"many2many:partssupply_serviceorders;"`
}

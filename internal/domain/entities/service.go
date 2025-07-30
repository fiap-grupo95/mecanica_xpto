package entities

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	ID                uint           `gorm:"primaryKey"`
	Name              string         `gorm:"size:100;not null"`
	Description       string         `gorm:"type:text"`
	Price             float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	AdditionalRepairs []AdditionalRepair
	ServiceOrders     []ServiceOrder `gorm:"many2many:service_serviceorders;"`
}

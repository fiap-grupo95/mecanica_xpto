package entities

import (
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID            uint           `gorm:"primaryKey"`
	Plate         string         `gorm:"size:10;not null"`
	CustomerID    uint           `gorm:"not null"`
	Customer      Customer       `gorm:"foreignKey:CustomerID"`
	Model         string         `gorm:"size:50;not null"`
	Year          string         `gorm:"size:4"`
	Brand         string         `gorm:"size:50;not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ServiceOrders []ServiceOrder
}

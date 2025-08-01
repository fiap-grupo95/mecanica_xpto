package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type VehicleModel struct {
	ID            uint           `gorm:"primaryKey"`
	Plate         string         `gorm:"size:10;not null"`
	CustomerID    uint           `gorm:"not null"`
	Customer      CustomerModel  `gorm:"foreignKey:CustomerID"`
	Model         string         `gorm:"size:50;not null"`
	Year          string         `gorm:"size:4"`
	Brand         string         `gorm:"size:50;not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ServiceOrders []ServiceOrderModel
}

func (v *VehicleModel) ToDomain() entities.Vehicle {
	return entities.Vehicle{
		ID:         v.ID,
		Plate:      v.Plate,
		CustomerID: v.CustomerID,
		Customer:   v.Customer.ToDomain(),
		Model:      v.Model,
		Year:       v.Year,
		Brand:      v.Brand,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
		DeletedAt: func() *time.Time {
			if v.DeletedAt.Valid {
				return &v.DeletedAt.Time
			}
			return nil
		}(),
		ServiceOrders: nil, // This will be populated by the repository layer
	}
}

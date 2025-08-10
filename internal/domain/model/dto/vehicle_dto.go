package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"time"

	"gorm.io/gorm"
)

type VehicleDTO struct {
	ID            uint              `gorm:"primaryKey"`
	Plate         string            `gorm:"size:10;not null"`
	CustomerID    uint              `gorm:"column:customer_id;not null"`
	Customer      *CustomerDTO      `gorm:"foreignKey:CustomerID"`
	Model         string            `gorm:"size:50;not null"`
	Year          string            `gorm:"size:4"`
	Brand         string            `gorm:"size:50;not null"`
	CreatedAt     time.Time         `gorm:"autoCreateTime"`
	UpdatedAt     *time.Time        `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt    `gorm:"index"`
	ServiceOrders []ServiceOrderDTO `gorm:"foreignKey:VehicleID;references:ID"`
}

func (v *VehicleDTO) ToDomain() *entities.Vehicle {
	var customer *entities.Customer

	if v.Customer != nil {
		customer = v.Customer.ToDomain()
	}

	return &entities.Vehicle{
		ID:         v.ID,
		Plate:      valueobject.ParsePlate(v.Plate),
		CustomerID: v.CustomerID,
		Customer:   customer,
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
	}
}

// TableName specifies the table name for VehicleDTO
func (v *VehicleDTO) TableName() string {
	return "tb_vehicle"
}

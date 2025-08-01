package repository

import "mecanica_xpto/internal/domain/entities"

// N:N relationship between PartsSupply and ServiceOrder
type PartsSupplyServiceOrderModel struct {
	PartsSupplyID  uint `gorm:"primaryKey"`
	ServiceOrderID uint `gorm:"primaryKey"`
	Quantity       int  `gorm:"not null;default:1"`
}

func (pssom *PartsSupplyServiceOrderModel) ToDomain() entities.PartsSupplyServiceOrder {
	return entities.PartsSupplyServiceOrder{
		PartsSupplyID:  pssom.PartsSupplyID,
		ServiceOrderID: pssom.ServiceOrderID,
		Quantity:       pssom.Quantity,
	}
}

package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
)

// 1:1 relationship between Customer and User
// 1:N relationship between Customer and Vehicle
// 1:N relationship between Customer and ServiceOrder
type CustomerDTO struct {
	ID            uint     `gorm:"primaryKey"`
	UserID        uint     `gorm:"unique;not null"`
	User          *UserDTO `gorm:"foreignKey:UserID;references:ID"`
	CpfCnpj       string   `gorm:"size:20;not null"`
	PhoneNumber   string   `gorm:"size:20;not null"`
	FullName      string   `gorm:"column:fullname;size:100;not null"`
	Vehicles      []VehicleDTO
	ServiceOrders []ServiceOrderDTO
}

func (cm *CustomerDTO) TableName() string {
	return "tb_customer"
}

func (cm *CustomerDTO) ToDomain() entities.Customer {
	var user *entities.User
	if cm.User != nil {
		u := cm.User.ToDomain()
		user = &u
	}

	return entities.Customer{
		ID:            cm.ID,
		UserID:        cm.UserID,
		User:          user,
		CpfCnpj:       valueobject.CpfCnpj(cm.CpfCnpj),
		PhoneNumber:   cm.PhoneNumber,
		FullName:      cm.FullName,
		Vehicles:      nil, // This will be populated later if needed
		ServiceOrders: nil, // This will be populated later if needed
	}
}

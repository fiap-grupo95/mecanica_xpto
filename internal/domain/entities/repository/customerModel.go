package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/entities/valueobject"
)

// 1:1 relationship between Customer and User
// 1:N relationship between Customer and Vehicle
// 1:N relationship between Customer and ServiceOrder
type CustomerModel struct {
	ID            uint       `gorm:"primaryKey"`
	UserID        uint       `gorm:"unique;not null"`
	User          *UserModel `gorm:"foreignKey:UserID;references:ID"`
	CPF_CNPJ      string     `gorm:"size:20;not null"`
	PhoneNumber   string     `gorm:"size:20;not null"`
	FullName      string     `gorm:"size:100;not null"`
	Vehicles      []VehicleModel
	ServiceOrders []ServiceOrderModel
}

func (cm *CustomerModel) ToDomain() entities.Customer {
	var user *entities.User
	if cm.User != nil {
		u := cm.User.ToDomain()
		user = &u
	}
	return entities.Customer{
		ID:            cm.ID,
		UserID:        cm.UserID,
		User:          user,
		CPF_CNPJ:      valueobject.ParseCPF_CNPJ(cm.CPF_CNPJ),
		PhoneNumber:   cm.PhoneNumber,
		FullName:      cm.FullName,
		Vehicles:      nil, // This will be populated later if needed
		ServiceOrders: nil, // This will be populated later if needed
	}
}

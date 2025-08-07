package dto

import (
	"log"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
)

// 1:1 relationship between Customer and User
// 1:N relationship between Customer and Vehicle
// 1:N relationship between Customer and ServiceOrder
type CustomerDTO struct {
	ID            uint              `gorm:"primaryKey"`
	UserID        uint              `gorm:"unique;not null"`
	User          *UserDTO          `gorm:"foreignKey:UserID;references:ID"`
	CpfCnpj       string            `gorm:"size:20;not null"`
	PhoneNumber   string            `gorm:"size:20;not null"`
	FullName      string            `gorm:"size:100;not null"`
	Vehicles      []VehicleDTO      `gorm:"foreignKey:CustomerID"`
	ServiceOrders []ServiceOrderDTO `gorm:"foreignKey:CustomerID"`
}

func (cm *CustomerDTO) TableName() string {
	return "customers"
}

func (cm *CustomerDTO) ToDomain() entities.Customer {
	var user *entities.User
	if cm.User != nil {
		u := cm.User.ToDomain()
		user = &u
	}

	cpfCnpj, err := valueobject.NewCpfCnpj(cm.CpfCnpj)
	if err != nil {
		log.Fatalf("Invalid CPF/CNPJ format: %v", err)
	}

	return entities.Customer{
		ID:            cm.ID,
		UserID:        cm.UserID,
		User:          user,
		CpfCnpj:       cpfCnpj,
		PhoneNumber:   cm.PhoneNumber,
		FullName:      cm.FullName,
		Vehicles:      nil, // This will be populated later if needed
		ServiceOrders: nil, // This will be populated later if needed
	}
}

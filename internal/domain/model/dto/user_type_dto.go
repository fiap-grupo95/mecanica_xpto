package dto

import (
	"mecanica_xpto/internal/domain/model/valueobject"
)

type UserTypeDTO struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gorm:"size:50;not null"`
	//Users []UserDTO `gorm:"foreignKey:UserTypeID"`
}

func (utm *UserTypeDTO) ToDomain() valueobject.UserType {
	return valueobject.ParseUserType(utm.Type)
}

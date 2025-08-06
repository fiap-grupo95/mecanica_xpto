package repository

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/valueobject"
)

type UserTypeDTO struct {
	ID    uint          `gorm:"primaryKey"`
	Type  string        `gorm:"size:50;not null"`
	Users []dto.UserDTO `gorm:"foreignKey:UserTypeID;references:ID"`
}

func (utm *UserTypeDTO) ToDomain() valueobject.UserType {
	return valueobject.ParseUserType(utm.Type)
}

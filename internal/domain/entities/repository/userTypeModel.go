package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"mecanica_xpto/internal/domain/entities/valueobject"
)

type UserTypeModel struct {
	ID    uint   `gorm:"primaryKey"`
	Type  string `gorm:"size:50;not null"`
	Users []UserModel
}

func (utm *UserTypeModel) ToDomain() entities.UserType {
	return entities.UserType{
		ID:   utm.ID,
		Type: valueobject.ParseUserType(utm.Type),
	}
}

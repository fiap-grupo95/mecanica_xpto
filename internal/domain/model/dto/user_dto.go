package dto

import (
	"mecanica_xpto/internal/domain/model/entities"
	"time"

	"gorm.io/gorm"
)

type UserDTO struct {
	ID         uint           `gorm:"primaryKey"`
	Email      string         `gorm:"size:100;not null;unique" json:"email" binding:"required,email"`
	Password   string         `gorm:"size:255;not null" json:"password" binding:"required"`
	UserTypeID uint           `gorm:"not null"`
	UserType   UserTypeDTO    `gorm:"foreignKey:UserTypeID"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Customer   *CustomerDTO   `gorm:"foreignKey:UserID;references:ID"`
}

func (m *UserDTO) ToDomain() entities.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	var customer *entities.Customer
	if m.Customer != nil {
		c := m.Customer.ToDomain()
		customer = &c
	}
	return entities.User{
		ID:        m.ID,
		Email:     m.Email,
		Password:  m.Password,
		UserType:  m.UserType.ToDomain(),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
		Customer:  customer,
	}
}

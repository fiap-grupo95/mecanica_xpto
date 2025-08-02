package repository

import (
	"mecanica_xpto/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID         uint           `gorm:"primaryKey"`
	Email      string         `gorm:"size:100;not null;unique"`
	Password   string         `gorm:"size:255;not null"`
	UserTypeID uint           `gorm:"not null"`
	UserType   UserTypeModel  `gorm:"foreignKey:UserTypeID"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Customer   *CustomerModel `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (m *UserModel) ToDomain() entities.User {
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
		UserType:  m.UserType.ToDomain().Type,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
		Customer:  customer,
	}
}

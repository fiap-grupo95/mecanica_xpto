package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey"`
	Email      string         `gorm:"size:100;not null;unique"`
	Password   string         `gorm:"size:255;not null"`
	UserTypeID uint           `gorm:"not null"`
	UserType   UserType       `gorm:"foreignKey:UserTypeID"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Customer   *Customer      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

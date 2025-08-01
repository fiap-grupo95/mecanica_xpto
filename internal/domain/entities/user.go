package entities

import (
	"time"
)

type User struct {
	ID         uint       `json:"id"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	UserTypeID uint       `json:"user_type_id"`
	UserType   UserType   `json:"user_type,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Customer   *Customer  `json:"customer,omitempty"`
}

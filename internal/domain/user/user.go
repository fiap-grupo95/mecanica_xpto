package user

import "errors"

// User represents the user entity
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // "-" means this field won't be included in JSON
}

// ErrUserNotFound is returned when a user cannot be found
var ErrUserNotFound = errors.New("user not found")

// Repository defines the interface for user data access
type Repository interface {
	GetByID(id string) (*User, error)
}

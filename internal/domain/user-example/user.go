package user_example

import "errors"

// User represents the user-example entity
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"` // "-" means this field won't be included in JSON
}

// ErrUserNotFound is returned when a user-example cannot be found
var ErrUserNotFound = errors.New("user-example not found")

// Repository defines the interface for user-example data access
type Repository interface {
	GetByID(id string) (*User, error)
}

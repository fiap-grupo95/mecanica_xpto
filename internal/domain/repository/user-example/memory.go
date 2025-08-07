package user_example

import (
	"sync"
)

// MemoryRepository implements user-example.Repository interface with in-memory storage
type MemoryRepository struct {
	users map[string]*User
	mu    sync.RWMutex
}

// NewMemoryRepository creates a new memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*User),
	}
}

// GetByID retrieves a user-example by their ID
func (r *MemoryRepository) GetByID(id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, exists := r.users[id]; exists {
		return u, nil
	}
	return nil, ErrUserNotFound
}

// AddUser adds a user-example to the repository (helper for testing)
func (r *MemoryRepository) AddUser(u *User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
}

// GetAllUsers retrieves all users from the repository (helper for testing)
func (r *MemoryRepository) GetAllUsers() map[string]*User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a new map to avoid exposing the internal map directly
	usersCopy := make(map[string]*User, len(r.users))
	for id, u := range r.users {
		usersCopy[id] = u
	}
	return usersCopy
}

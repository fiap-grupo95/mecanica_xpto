package repository

import (
	"mecanica_xpto/internal/domain/user-example"
	"sync"
)

// MemoryRepository implements user-example.Repository interface with in-memory storage
type MemoryRepository struct {
	users map[string]*user_example.User
	mu    sync.RWMutex
}

// NewMemoryRepository creates a new memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*user_example.User),
	}
}

// GetByID retrieves a user-example by their ID
func (r *MemoryRepository) GetByID(id string) (*user_example.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, exists := r.users[id]; exists {
		return u, nil
	}
	return nil, user_example.ErrUserNotFound
}

// AddUser adds a user-example to the repository (helper for testing)
func (r *MemoryRepository) AddUser(u *user_example.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
}

// GetAllUsers retrieves all users from the repository (helper for testing)
func (r *MemoryRepository) GetAllUsers() map[string]*user_example.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a new map to avoid exposing the internal map directly
	usersCopy := make(map[string]*user_example.User, len(r.users))
	for id, u := range r.users {
		usersCopy[id] = u
	}
	return usersCopy
}

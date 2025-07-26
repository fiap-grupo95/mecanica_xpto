package repository

import (
	"mecanica_xpto/internal/domain/user"
	"sync"
)

// MemoryRepository implements user.Repository interface with in-memory storage
type MemoryRepository struct {
	users map[string]*user.User
	mu    sync.RWMutex
}

// NewMemoryRepository creates a new memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*user.User),
	}
}

// GetByID retrieves a user by their ID
func (r *MemoryRepository) GetByID(id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, exists := r.users[id]; exists {
		return u, nil
	}
	return nil, user.ErrUserNotFound
}

// AddUser adds a user to the repository (helper for testing)
func (r *MemoryRepository) AddUser(u *user.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
}

// GetAllUsers retrieves all users from the repository (helper for testing)
func (r *MemoryRepository) GetAllUsers() map[string]*user.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a new map to avoid exposing the internal map directly
	usersCopy := make(map[string]*user.User, len(r.users))
	for id, u := range r.users {
		usersCopy[id] = u
	}
	return usersCopy
}

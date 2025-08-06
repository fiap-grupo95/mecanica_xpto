package customers

import (
	"mecanica_xpto/internal/domain/model/dto"
)

// ICustomerRepository defines the interface for customers data access
type ICustomerRepository interface {
	GetByID(id uint) (*dto.CustomerDTO, error)
	Create(customer *dto.CustomerDTO) error
	Update(customer *dto.CustomerDTO) error
	Delete(id uint) error
	List() ([]dto.CustomerDTO, error)
}

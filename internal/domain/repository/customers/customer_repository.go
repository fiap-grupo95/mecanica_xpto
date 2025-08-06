package customers

import (
	"gorm.io/gorm"
	"mecanica_xpto/internal/domain/model/dto"
)

// CustomerRepository implements ICustomerRepository interface
type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) ICustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *dto.CustomerDTO) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) GetByID(id uint) (*dto.CustomerDTO, error) {
	var customer dto.CustomerDTO
	err := r.db.Preload("User").First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Update(customer *dto.CustomerDTO) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(id uint) error {
	return r.db.Delete(&dto.CustomerDTO{}, id).Error
}

func (r *CustomerRepository) List() ([]dto.CustomerDTO, error) {
	var customers []dto.CustomerDTO
	err := r.db.Preload("User").Find(&customers).Error
	return customers, err
}

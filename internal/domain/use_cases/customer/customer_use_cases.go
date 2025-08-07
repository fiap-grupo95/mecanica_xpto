package use_cases

import (
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	"mecanica_xpto/internal/domain/repository/customers"
)

// ICustomerUseCase defines the interface for customers use cases
type ICustomerUseCase interface {
	GetById(id uint) (*entities.Customer, error)
	GetByDocument(CpfCnpj string) (*entities.Customer, error)
	CreateCustomer(customer *entities.Customer) error
	UpdateCustomer(id uint, customer *entities.Customer) error
	DeleteCustomer(id uint) error
	ListCustomer() ([]entities.Customer, error)
}
type CustomerUseCase struct {
	customerRepo repository.ICustomerRepository
	userRepo     repository.IUserRepository
}

func NewCustomerUseCase(customerRepo repository.ICustomerRepository, userRepo repository.IUserRepository) ICustomerUseCase {
	return &CustomerUseCase{customerRepo: customerRepo, userRepo: userRepo}
}

func (uc *CustomerUseCase) GetById(id uint) (*entities.Customer, error) {
	customerDTO, err := uc.customerRepo.GetByID(id)

	if err != nil {
		return nil, err
	}

	return customerDTO.ToDomain(), nil
}

func (uc *CustomerUseCase) GetByDocument(CpfCnpj string) (*entities.Customer, error) {
	customerDTO, err := uc.customerRepo.GetByDocument(CpfCnpj)

	if err != nil {
		return nil, err
	}

	return customerDTO.ToDomain(), nil
}

func (uc *CustomerUseCase) CreateCustomer(customer *entities.Customer) error {
	userDTO := dto.UserDTO{
		Email:    customer.Email,
		UserType: valueobject.ParseUserType("admin"),
	}
	customerDTO := dto.CustomerDTO{
		User:        &userDTO,
		CpfCnpj:     customer.CpfCnpj.String(),
		PhoneNumber: customer.PhoneNumber,
		FullName:    customer.FullName,
	}

	return uc.customerRepo.Create(&customerDTO)
}

func (uc *CustomerUseCase) UpdateCustomer(id uint, customer *entities.Customer) error {
	existingDTO, err := uc.customerRepo.GetByID(id)
	if err != nil {
		return err
	}
	if customer.FullName != "" {
		existingDTO.FullName = customer.FullName
	}

	if customer.PhoneNumber != "" {
		existingDTO.PhoneNumber = customer.PhoneNumber
	}

	return uc.customerRepo.Update(existingDTO)
}

func (uc *CustomerUseCase) DeleteCustomer(id uint) error {
	return uc.customerRepo.Delete(id)
}

func (uc *CustomerUseCase) ListCustomer() ([]entities.Customer, error) {
	dtos, err := uc.customerRepo.List()
	if err != nil {
		return nil, err
	}
	customers := make([]entities.Customer, 0, len(dtos))
	for _, dto := range dtos {
		customers = append(customers, *dto.ToDomain())
	}
	return customers, nil
}

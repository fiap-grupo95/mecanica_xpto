package usecase

import (
	"errors"
	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/model/valueobject"
	customerRepo "mecanica_xpto/internal/domain/repository/customers"
	"mecanica_xpto/internal/domain/repository/users"
)

var (
	ErrGeneric               = errors.New("unknown error")
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrInvalidDocumentFormat = errors.New("invalid document format")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidCustomerID     = errors.New("invalid customer ID")
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
	customerRepo customerRepo.ICustomerRepository
	userRepo     users.IUserRepository
}

func NewCustomerUseCase(customerRepo customerRepo.ICustomerRepository, userRepo users.IUserRepository) ICustomerUseCase {
	return &CustomerUseCase{customerRepo: customerRepo, userRepo: userRepo}
}

func (uc *CustomerUseCase) GetById(id uint) (*entities.Customer, error) {
	customerDTO, err := uc.customerRepo.GetByID(id)

	if err != nil {
		return nil, ErrGeneric
	}
	if customerDTO == nil {
		return nil, ErrCustomerNotFound
	}

	return customerDTO.ToDomain(), nil
}

func (uc *CustomerUseCase) GetByDocument(CpfCnpj string) (*entities.Customer, error) {
	customerDTO, err := uc.customerRepo.GetByDocument(CpfCnpj)

	if err != nil {
		return nil, ErrGeneric
	}
	if customerDTO == nil {
		return nil, ErrCustomerNotFound
	}

	return customerDTO.ToDomain(), nil
}

func (uc *CustomerUseCase) CreateCustomer(customer *entities.Customer) error {
	if e := customer.CpfCnpj.IsValid(); e != nil {
		return ErrInvalidDocumentFormat
	}
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

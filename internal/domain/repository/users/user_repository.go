package users

import (
	"mecanica_xpto/internal/domain/model/dto"

	"gorm.io/gorm"
)

// IUserRepository defines the interface for Users data access
type IUserRepository interface {
	GetByID(id uint) (*dto.UserDTO, error)
	GetByEmail(email string) (*dto.UserDTO, error)
	Create(User *dto.UserDTO) error
	Update(User *dto.UserDTO) error
	Delete(id uint) error
	List() ([]dto.UserDTO, error)
}

// UserRepository implements IUserRepository interface
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(User *dto.UserDTO) error {
	return r.db.Create(User).Error
}

func (r *UserRepository) GetByID(id uint) (*dto.UserDTO, error) {
	var User dto.UserDTO
	err := r.db.Preload("User").First(&User, id).Error
	if err != nil {
		return nil, err
	}
	return &User, nil
}

func (r *UserRepository) GetByEmail(email string) (*dto.UserDTO, error) {
	var user dto.UserDTO
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(User *dto.UserDTO) error {
	return r.db.Save(User).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&dto.UserDTO{}, id).Error
}

func (r *UserRepository) List() ([]dto.UserDTO, error) {
	var Users []dto.UserDTO
	err := r.db.Preload("User").Find(&Users).Error
	return Users, err
}

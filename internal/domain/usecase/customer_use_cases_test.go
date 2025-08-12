package usecase_test

import (
	"errors"
	use_cases "mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/internal/domain/usecase/mocks"
	"testing"

	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	customer := &entities.Customer{FullName: "Updated", CpfCnpj: "123", PhoneNumber: "999"}
	dtoOrig := &dto.CustomerDTO{ID: 1, FullName: "Old", CpfCnpj: "123", PhoneNumber: "888", User: &dto.UserDTO{ID: 2, Email: "old@mail.com"}}

	mockRepo.EXPECT().GetByID(uint(1)).Return(dtoOrig, nil)
	mockRepo.EXPECT().Update(dtoOrig).Return(nil)

	err := uc.UpdateCustomer(1, customer)
	assert.NoError(t, err)
}

func TestUpdateCustomer_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().GetByID(uint(1)).Return(nil, errors.New("not found"))

	err := uc.UpdateCustomer(1, &entities.Customer{})
	assert.Error(t, err)
}

func TestDeleteCustomer_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().Delete(uint(1)).Return(nil)

	err := uc.DeleteCustomer(1)
	assert.NoError(t, err)
}

func TestDeleteCustomer_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().Delete(uint(1)).Return(errors.New("fail"))

	err := uc.DeleteCustomer(1)
	assert.Error(t, err)
}

func TestListCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockICustomerRepository(ctrl)
	uc := use_cases.NewCustomerUseCase(mockRepo, nil)

	mockRepo.EXPECT().List().Return(nil, errors.New("fail"))

	customers, err := uc.ListCustomer()
	assert.Error(t, err)
	assert.Nil(t, customers)
}

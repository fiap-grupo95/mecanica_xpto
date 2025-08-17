package usecase_test

import (
	"context"
	"errors"
	"testing"

	"mecanica_xpto/internal/domain/model/dto"
	"mecanica_xpto/internal/domain/model/entities"
	"mecanica_xpto/internal/domain/usecase"
	"mecanica_xpto/internal/domain/usecase/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdditionalRepair_ServiceOrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdditionalRepairRepository(ctrl)
	mockRepoOS := mocks.NewMockIServiceOrderRepository(ctrl)
	mockServiceRepo := mocks.NewMockIServiceRepo(ctrl)
	mockPartsSupplyRepo := mocks.NewMockIPartsSupplyRepo(ctrl)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{ServiceOrderID: 2}

	mockRepoOS.EXPECT().GetByID(uint(2)).Return(nil, errors.New("not found"))

	err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
}

func TestCreateAdditionalRepair_ServiceNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdditionalRepairRepository(ctrl)
	mockRepoOS := mocks.NewMockIServiceOrderRepository(ctrl)
	mockServiceRepo := mocks.NewMockIServiceRepo(ctrl)
	mockPartsSupplyRepo := mocks.NewMockIPartsSupplyRepo(ctrl)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: 1,
		Services:       []entities.Service{{ID: 99}},
	}

	mockRepoOS.EXPECT().GetByID(uint(1)).Return(&dto.ServiceOrderDTO{ID: 1}, nil)
	mockServiceRepo.EXPECT().GetByID(gomock.Any(), uint(99)).Return(entities.Service{}, nil)

	err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
}

func TestCreateAdditionalRepair_PartsSupplyNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAdditionalRepairRepository(ctrl)
	mockRepoOS := mocks.NewMockIServiceOrderRepository(ctrl)
	mockServiceRepo := mocks.NewMockIServiceRepo(ctrl)
	mockPartsSupplyRepo := mocks.NewMockIPartsSupplyRepo(ctrl)

	uc := usecase.NewSOAdditionalRepairUseCase(mockRepo, mockRepoOS, mockServiceRepo, mockPartsSupplyRepo)

	adr := entities.AdditionalRepair{
		ServiceOrderID: 1,
		PartsSupplies:  []entities.PartsSupply{{ID: 88}},
	}

	mockRepoOS.EXPECT().GetByID(uint(1)).Return(&dto.ServiceOrderDTO{ID: 1}, nil)
	mockPartsSupplyRepo.EXPECT().GetByID(gomock.Any(), uint(88)).Return(entities.PartsSupply{}, nil)

	err := uc.CreateAdditionalRepair(context.Background(), adr)
	assert.Error(t, err)
}

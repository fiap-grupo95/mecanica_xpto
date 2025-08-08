package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"mecanica_xpto/internal/domain/mocks"
	"mecanica_xpto/internal/domain/model/entities"

	"go.uber.org/mock/gomock"
)

func TestGetServiceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()

	service := entities.Service{ID: 1, Name: "Troca de óleo"}
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(service, nil)
	result, err := uc.GetServiceByID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, service) {
		t.Errorf("expected %v, got %v", service, result)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.Service{}, nil)
	_, err = uc.GetServiceByID(ctx, 2)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.Service{}, errors.New("fail"))
	_, err = uc.GetServiceByID(ctx, 3)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCreateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()
	service := &entities.Service{ID: 1, Name: "Troca de óleo"}

	mockRepo.EXPECT().GetByName(ctx, service.Name).Return(entities.Service{}, nil).AnyTimes()
	mockRepo.EXPECT().Create(ctx, service).Return(*service, nil)
	result, err := uc.CreateService(ctx, service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, *service) {
		t.Errorf("expected %v, got %v", *service, result)
	}

	mockRepo.EXPECT().GetByName(ctx, service.Name).Return(entities.Service{}, nil).AnyTimes()
	mockRepo.EXPECT().Create(ctx, service).Return(entities.Service{}, errors.New("fail"))
	_, err = uc.CreateService(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()
	service := &entities.Service{ID: 1, Name: "Troca de óleo"}

	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(*service, nil)
	mockRepo.EXPECT().Update(ctx, service).Return(nil)
	err := uc.UpdateService(ctx, service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.Service{}, nil)
	service.ID = 2
	err = uc.UpdateService(ctx, service)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(*service, nil)
	mockRepo.EXPECT().Update(ctx, service).Return(errors.New("fail"))
	service.ID = 3
	err = uc.UpdateService(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().Delete(ctx, uint(1)).Return(nil)
	err := uc.DeleteService(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().Delete(ctx, uint(2)).Return(errors.New("fail"))
	err = uc.DeleteService(ctx, 2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	uc := NewServiceUseCase(mockRepo)
	ctx := context.Background()
	services := []entities.Service{{ID: 1, Name: "Troca de óleo"}, {ID: 2, Name: "Alinhamento"}}

	mockRepo.EXPECT().List(ctx).Return(services, nil)
	result, err := uc.ListServices(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, services) {
		t.Errorf("expected %v, got %v", services, result)
	}

	mockRepo.EXPECT().List(ctx).Return(nil, errors.New("fail"))
	_, err = uc.ListServices(ctx)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

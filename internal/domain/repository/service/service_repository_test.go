package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"mecanica_xpto/internal/domain/mocks"
	"mecanica_xpto/internal/domain/model/entities"

	"go.uber.org/mock/gomock"
)

func TestIServiceRepo_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	ctx := context.Background()
	service := &entities.Service{ID: 1, Name: "Troca de óleo", Description: "Troca de óleo sintético", Price: 100}

	mockRepo.EXPECT().Create(ctx, service).Return(*service, nil)
	result, err := mockRepo.Create(ctx, service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, *service) {
		t.Errorf("expected %v, got %v", *service, result)
	}

	mockRepo.EXPECT().Create(ctx, service).Return(entities.Service{}, errors.New("fail"))
	_, err = mockRepo.Create(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIServiceRepo_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	ctx := context.Background()
	service := entities.Service{ID: 1, Name: "Troca de óleo"}

	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(service, nil)
	result, err := mockRepo.GetByID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, service) {
		t.Errorf("expected %v, got %v", service, result)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.Service{}, errors.New("not found"))
	_, err = mockRepo.GetByID(ctx, 2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIServiceRepo_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	ctx := context.Background()
	service := &entities.Service{ID: 1, Name: "Troca de óleo"}

	mockRepo.EXPECT().Update(ctx, service).Return(nil)
	err := mockRepo.Update(ctx, service)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().Update(ctx, service).Return(errors.New("fail"))
	err = mockRepo.Update(ctx, service)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIServiceRepo_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	ctx := context.Background()

	mockRepo.EXPECT().Delete(ctx, uint(1)).Return(nil)
	err := mockRepo.Delete(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().Delete(ctx, uint(2)).Return(errors.New("fail"))
	err = mockRepo.Delete(ctx, 2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIServiceRepo_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIServiceRepo(ctrl)
	ctx := context.Background()
	services := []entities.Service{{ID: 1, Name: "Troca de óleo"}, {ID: 2, Name: "Alinhamento"}}

	mockRepo.EXPECT().List(ctx).Return(services, nil)
	result, err := mockRepo.List(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, services) {
		t.Errorf("expected %v, got %v", services, result)
	}

	mockRepo.EXPECT().List(ctx).Return(nil, errors.New("fail"))
	_, err = mockRepo.List(ctx)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

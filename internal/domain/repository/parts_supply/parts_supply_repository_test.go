package parts_supply

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"mecanica_xpto/internal/domain/mocks"
	"mecanica_xpto/internal/domain/model/entities"

	"go.uber.org/mock/gomock"
)

func TestIPartsSupplyRepo_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	ctx := context.Background()
	ps := &entities.PartsSupply{ID: 1, Name: "Filtro de óleo", Description: "Filtro de óleo sintético", Price: 30, QuantityTotal: 10, QuantityReserve: 2}

	mockRepo.EXPECT().Create(ctx, ps).Return(*ps, nil)
	result, err := mockRepo.Create(ctx, ps)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, *ps) {
		t.Errorf("expected %v, got %v", *ps, result)
	}

	mockRepo.EXPECT().Create(ctx, ps).Return(entities.PartsSupply{}, errors.New("fail"))
	_, err = mockRepo.Create(ctx, ps)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIPartsSupplyRepo_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	ctx := context.Background()
	ps := entities.PartsSupply{ID: 1, Name: "Filtro de óleo"}

	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(ps, nil)
	result, err := mockRepo.GetByID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, ps) {
		t.Errorf("expected %v, got %v", ps, result)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.PartsSupply{}, errors.New("not found"))
	_, err = mockRepo.GetByID(ctx, 2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIPartsSupplyRepo_GetByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	ctx := context.Background()
	ps := entities.PartsSupply{ID: 1, Name: "Filtro de óleo"}

	mockRepo.EXPECT().GetByName(ctx, "Filtro de óleo").Return(ps, nil)
	result, err := mockRepo.GetByName(ctx, "Filtro de óleo")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, ps) {
		t.Errorf("expected %v, got %v", ps, result)
	}

	mockRepo.EXPECT().GetByName(ctx, "Outro").Return(entities.PartsSupply{}, errors.New("not found"))
	_, err = mockRepo.GetByName(ctx, "Outro")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIPartsSupplyRepo_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	ctx := context.Background()
	ps := &entities.PartsSupply{ID: 1, Name: "Filtro de óleo"}

	mockRepo.EXPECT().Update(ctx, ps).Return(nil)
	err := mockRepo.Update(ctx, ps)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().Update(ctx, ps).Return(errors.New("fail"))
	err = mockRepo.Update(ctx, ps)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIPartsSupplyRepo_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
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

func TestIPartsSupplyRepo_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	ctx := context.Background()
	parts := []entities.PartsSupply{{ID: 1, Name: "Filtro de óleo"}, {ID: 2, Name: "Pastilha de freio"}}

	mockRepo.EXPECT().List(ctx).Return(parts, nil)
	result, err := mockRepo.List(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, parts) {
		t.Errorf("expected %v, got %v", parts, result)
	}

	mockRepo.EXPECT().List(ctx).Return(nil, errors.New("fail"))
	_, err = mockRepo.List(ctx)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

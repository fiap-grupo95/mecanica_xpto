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

func TestGetPartsSupplyByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()

	ps := entities.PartsSupply{ID: 1, Name: "Filtro"}
	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(ps, nil)
	result, err := uc.GetPartsSupplyByID(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, ps) {
		t.Errorf("expected %v, got %v", ps, result)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.PartsSupply{}, nil)
	_, err = uc.GetPartsSupplyByID(ctx, 2)
	if !errors.Is(err, ErrPartsSupplyNotFound) {
		t.Errorf("expected ErrPartsSupplyNotFound, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.PartsSupply{}, errors.New("fail"))
	_, err = uc.GetPartsSupplyByID(ctx, 3)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCreatePartsSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()
	ps := &entities.PartsSupply{ID: 1, Name: "Filtro"}

	mockRepo.EXPECT().GetByName(ctx, "Filtro").Return(entities.PartsSupply{ID: 1, Name: "Filtro"}, nil)
	_, err := uc.CreatePartsSupply(ctx, ps)
	if !errors.Is(err, ErrPartsSupplyAlreadyExists) {
		t.Errorf("expected ErrPartsSupplyAlreadyExists, got %v", err)
	}

	mockRepo.EXPECT().GetByName(ctx, "Novo").Return(entities.PartsSupply{}, errors.New("not found"))
	mockRepo.EXPECT().Create(ctx, ps).Return(*ps, nil)
	ps.Name = "Novo"
	_, err = uc.CreatePartsSupply(ctx, ps)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().GetByName(ctx, "Falha").Return(entities.PartsSupply{}, errors.New("not found"))
	mockRepo.EXPECT().Create(ctx, ps).Return(entities.PartsSupply{}, errors.New("fail"))
	ps.Name = "Falha"
	_, err = uc.CreatePartsSupply(ctx, ps)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestUpdatePartsSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()
	ps := &entities.PartsSupply{ID: 1, Name: "Filtro"}

	mockRepo.EXPECT().GetByID(ctx, uint(1)).Return(entities.PartsSupply{ID: 1, Name: "Filtro"}, nil)
	mockRepo.EXPECT().Update(ctx, ps).Return(nil)
	err := uc.UpdatePartsSupply(ctx, ps)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(2)).Return(entities.PartsSupply{}, nil)
	ps.ID = 2
	err = uc.UpdatePartsSupply(ctx, ps)
	if err == nil || err.Error() != "parts supply not found" {
		t.Errorf("expected not found error, got %v", err)
	}

	mockRepo.EXPECT().GetByID(ctx, uint(3)).Return(entities.PartsSupply{ID: 3, Name: "Falha"}, nil)
	mockRepo.EXPECT().Update(ctx, ps).Return(errors.New("fail"))
	ps.ID = 3
	err = uc.UpdatePartsSupply(ctx, ps)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDeletePartsSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()

	mockRepo.EXPECT().Delete(ctx, uint(1)).Return(nil)
	err := uc.DeletePartsSupply(ctx, 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRepo.EXPECT().Delete(ctx, uint(2)).Return(errors.New("fail"))
	err = uc.DeletePartsSupply(ctx, 2)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestListPartsSupplies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockIPartsSupplyRepo(ctrl)
	uc := NewPartsSupplyUseCase(mockRepo)
	ctx := context.Background()
	parts := []entities.PartsSupply{{ID: 1, Name: "Filtro"}, {ID: 2, Name: "Pastilha"}}

	mockRepo.EXPECT().List(ctx).Return(parts, nil)
	result, err := uc.ListPartsSupplies(ctx)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, parts) {
		t.Errorf("expected %v, got %v", parts, result)
	}

	mockRepo.EXPECT().List(ctx).Return(nil, errors.New("fail"))
	_, err = uc.ListPartsSupplies(ctx)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

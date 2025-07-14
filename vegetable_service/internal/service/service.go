package service

import (
	"context"

	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/entity"
)

type VegetableService interface {
	CreateVegetable(ctx context.Context, in dto.CreateVegetableParams) (string, error)
	DeleteVegetable(ctx context.Context, vegetableID string) error
	GetVegetable(ctx context.Context, vegetableID string) (entity.Vegetable, error)
	ListVegetables(ctx context.Context, in dto.ListVegetablesParams) ([]entity.Vegetable, error)
	UpdateVegetable(ctx context.Context, in dto.VegetableUpdateOperation) error
	// TODO: updated resource must be returned
	BatchUpdateVegetables(ctx context.Context, in dto.VegetablesBatchUpdateOperation) error
}

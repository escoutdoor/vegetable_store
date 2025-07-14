package client

import (
	"context"

	vegetable_client "github.com/escoutdoor/vegetable_store/order_service/internal/client/vegetable"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
)

type VegetableClient interface {
	GetVegetable(ctx context.Context, vegetableID string) (entity.Vegetable, error)
	ListVegetables(ctx context.Context, in vegetable_client.ListVegetablesParams) (map[string]*entity.Vegetable, error)
	BatchUpdateVegetablesWeight(ctx context.Context, in []vegetable_client.UpdateVegetableWeightParams) error
}

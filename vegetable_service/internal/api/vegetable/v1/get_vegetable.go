package v1

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/converter"
)

func (i *Implementation) GetVegetable(ctx context.Context, req *vegetablev1.GetVegetableRequest) (*vegetablev1.GetVegetableResponse, error) {
	vegetable, err := i.vegetableService.GetVegetable(ctx, req.GetVegetableId())
	if err != nil {
		return nil, err
	}

	return &vegetablev1.GetVegetableResponse{Vegetable: converter.VegetableToProto(vegetable)}, nil
}

package v1

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
)

func (i *Implementation) DeleteVegetable(ctx context.Context, req *vegetablev1.DeleteVegetableRequest) (*vegetablev1.DeleteVegetableResponse, error) {
	err := i.vegetableService.DeleteVegetable(ctx, req.GetVegetableId())
	if err != nil {
		return nil, err
	}

	return &vegetablev1.DeleteVegetableResponse{}, nil
}

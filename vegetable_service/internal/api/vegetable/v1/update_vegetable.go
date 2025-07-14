package v1

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/converter"
)

func (i *Implementation) UpdateVegetable(ctx context.Context, req *vegetablev1.UpdateVegetableRequest) (*vegetablev1.UpdateVegetableResponse, error) {
	if err := validateUpdateVegetableRequest(req); err != nil {
		return nil, err
	}

	in, err := converter.ProtoUpdateVegetableRequestToVegetableUpdateOperation(req)
	if err != nil {
		return nil, err
	}

	if err := i.vegetableService.UpdateVegetable(ctx, in); err != nil {
		return nil, err
	}

	return &vegetablev1.UpdateVegetableResponse{}, nil
}

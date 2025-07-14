package v1

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/converter"
)

func (i *Implementation) CreateVegetable(ctx context.Context, req *vegetablev1.CreateVegetableRequest) (*vegetablev1.CreateVegetableResponse, error) {
	id, err := i.vegetableService.CreateVegetable(ctx, converter.ProtoCreateVegetableRequestToCreateVegetableParams(req))
	if err != nil {
		return nil, err
	}

	return &vegetablev1.CreateVegetableResponse{VegetableId: id}, nil
}

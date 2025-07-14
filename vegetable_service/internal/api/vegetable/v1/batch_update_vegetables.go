package v1

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/converter"
)

func (i *Implementation) BatchUpdateVegetables(ctx context.Context, req *vegetablev1.BatchUpdateVegetablesRequest) (*vegetablev1.BatchUpdateVegetablesResponse, error) {
	if err := validateBatchUpdateVegetablesRequest(req); err != nil {
		return nil, err
	}

	in, err := converter.ProtoBatchUpdateVegetablesRequestToVegetablesBatchUpdateOperation(req)
	if err != nil {
		return nil, err
	}

	if err := i.vegetableService.BatchUpdateVegetables(ctx, in); err != nil {
		return nil, err
	}

	return &vegetablev1.BatchUpdateVegetablesResponse{}, nil
}

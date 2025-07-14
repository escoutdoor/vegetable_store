package v1

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/converter"
)

func (i *Implementation) ListVegetables(ctx context.Context, req *vegetablev1.ListVegetablesRequest) (*vegetablev1.ListVegetablesResponse, error) {
	list, err := i.vegetableService.ListVegetables(ctx, converter.ProtoListVegetablesRequestToListVegetablesParams(req))
	if err != nil {
		return nil, err
	}

	return &vegetablev1.ListVegetablesResponse{
		Vegetables: converter.VegetableListToProto(list),
		TotalSize:  int64(len((list))),
	}, nil
}

package vegetable

import (
	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
)

func protoVegetableToVegetable(vegetable *vegetablev1.Vegetable) *entity.Vegetable {
	return &entity.Vegetable{
		ID:              vegetable.Id,
		Name:            vegetable.Name,
		Weight:          vegetable.Weight,
		Price:           vegetable.Price,
		DiscountedPrice: vegetable.DiscountedPrice,
	}
}

func protoVegetablesToVegetables(vegetables []*vegetablev1.Vegetable) []*entity.Vegetable {
	list := make([]*entity.Vegetable, 0, len(vegetables))
	for _, v := range vegetables {
		list = append(list, protoVegetableToVegetable(v))
	}

	return list
}

func listVegetablesParamsToProtoListVegetablesRequest(params ListVegetablesParams) *vegetablev1.ListVegetablesRequest {
	return &vegetablev1.ListVegetablesRequest{
		Limit:        int64(params.Limit),
		Offset:       int64(params.Offset),
		VegetableIds: params.VegetableIDs,
	}
}

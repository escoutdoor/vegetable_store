package converter

import (
	"fmt"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/entity"

	"github.com/gojaguar/jaguar/strings"
	fieldmask_utils "github.com/mennanov/fieldmask-utils"
)

func ProtoCreateVegetableRequestToCreateVegetableParams(req *vegetablev1.CreateVegetableRequest) dto.CreateVegetableParams {
	return dto.CreateVegetableParams{
		Name:            req.Name,
		Weight:          req.Weight,
		Price:           req.Price,
		DiscountedPrice: req.DiscountedPrice,
	}
}

func ProtoUpdateVegetableRequestToVegetableUpdateOperation(req *vegetablev1.UpdateVegetableRequest) (dto.VegetableUpdateOperation, error) {
	update := &vegetablev1.VegetableUpdateOperation{}
	out := dto.VegetableUpdateOperation{ID: req.Update.GetVegetableId()}

	mask, err := fieldmask_utils.MaskFromProtoFieldMask(req.GetUpdateMask(), strings.PascalCase)
	if err != nil {
		return dto.VegetableUpdateOperation{}, errwrap.Wrap("create mask from the field mask", err)
	}

	if err := fieldmask_utils.StructToStruct(mask, req.Update, update); err != nil {
		return dto.VegetableUpdateOperation{}, errwrap.Wrap("copy struct to struct", err)
	}

	for _, p := range req.GetUpdateMask().GetPaths() {
		switch p {
		case "name":
			out.Name = &update.Name
		case "weight":
			out.Weight = &update.Weight
		case "price":
			out.Price = &update.Price
		case "discounted_price":
			out.DiscountedPrice = &update.DiscountedPrice
		}
	}

	return out, nil
}

func ProtoListVegetablesRequestToListVegetablesParams(req *vegetablev1.ListVegetablesRequest) dto.ListVegetablesParams {
	return dto.ListVegetablesParams{
		Limit:        int(req.Limit),
		Offset:       int(req.Offset),
		VegetableIDs: req.VegetableIds,
	}
}

func ProtoBatchUpdateVegetablesRequestToVegetablesBatchUpdateOperation(req *vegetablev1.BatchUpdateVegetablesRequest) (dto.VegetablesBatchUpdateOperation, error) {
	requests := make([]dto.VegetableUpdateOperation, 0, len(req.Requests))
	for i, r := range req.Requests {
		request, err := ProtoUpdateVegetableRequestToVegetableUpdateOperation(r)
		if err != nil {
			return dto.VegetablesBatchUpdateOperation{}, fmt.Errorf("convert %d request: %w", i, err)
		}
		requests = append(requests, request)
	}

	return dto.VegetablesBatchUpdateOperation{Requests: requests}, nil
}

func VegetableToProto(vegetable entity.Vegetable) *vegetablev1.Vegetable {
	return &vegetablev1.Vegetable{
		Id:              vegetable.ID,
		Name:            vegetable.Name,
		Weight:          vegetable.Weight,
		Price:           vegetable.Price,
		DiscountedPrice: vegetable.DiscountedPrice,
	}
}

func VegetableListToProto(vegetableList []entity.Vegetable) []*vegetablev1.Vegetable {
	list := make([]*vegetablev1.Vegetable, 0, len(vegetableList))
	for _, v := range vegetableList {
		list = append(list, VegetableToProto(v))
	}

	return list
}

package vegetable

import (
	"context"

	vegetablev1 "github.com/escoutdoor/vegetable_store/common/pkg/api/vegetable/v1"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type client struct {
	vegetableGrpcClient vegetablev1.VegetableServiceClient
}

func NewVegetableClient(vegetableGrpcClient vegetablev1.VegetableServiceClient) *client {
	return &client{
		vegetableGrpcClient: vegetableGrpcClient,
	}
}

type UpdateVegetableWeightParams struct {
	VegetableID string
	Weight      float32
}

type ListVegetablesParams struct {
	Limit        int
	Offset       int
	VegetableIDs []string
}

func (c *client) GetVegetable(ctx context.Context, vegetableID string) (entity.Vegetable, error) {
	req := &vegetablev1.GetVegetableRequest{VegetableId: vegetableID}
	resp, err := c.vegetableGrpcClient.GetVegetable(ctx, req)
	if err != nil {
		return entity.Vegetable{}, errwrap.Wrap("get vegetable using vegetable grpc client", err)
	}

	vegetable := protoVegetableToVegetable(resp.Vegetable)
	return *vegetable, nil
}

func (c *client) ListVegetables(ctx context.Context, in ListVegetablesParams) (map[string]*entity.Vegetable, error) {
	resp, err := c.vegetableGrpcClient.ListVegetables(ctx, listVegetablesParamsToProtoListVegetablesRequest(in))
	if err != nil {
		return nil, errwrap.Wrap("get the list of vegetables from vegetable grpc client", err)
	}
	vegetables := protoVegetablesToVegetables(resp.Vegetables)

	vegetablesMap := make(map[string]*entity.Vegetable, resp.TotalSize)
	for _, v := range vegetables {
		vegetablesMap[v.ID] = v
	}

	return vegetablesMap, nil
}

func (c *client) BatchUpdateVegetablesWeight(ctx context.Context, in []UpdateVegetableWeightParams) error {
	requests := make([]*vegetablev1.UpdateVegetableRequest, 0, len(in))

	for _, p := range in {
		requests = append(requests, &vegetablev1.UpdateVegetableRequest{
			Update:     &vegetablev1.VegetableUpdateOperation{VegetableId: p.VegetableID, Weight: p.Weight},
			UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"weight"}},
		})
	}

	req := &vegetablev1.BatchUpdateVegetablesRequest{Requests: requests}
	if _, err := c.vegetableGrpcClient.BatchUpdateVegetables(ctx, req); err != nil {
		return errwrap.Wrap("batch update vegetables using vegetable grpc client", err)
	}

	return nil
}

package v1

import (
	"context"

	orderv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/converter"
)

func (i *Implementation) ListOrders(ctx context.Context, req *orderv1.ListOrdersRequest) (*orderv1.ListOrdersResponse, error) {
	orders, err := i.orderService.ListOrders(ctx, converter.ProtoListOrdersRequestToListOrdersParams(req))
	if err != nil {
		return nil, err
	}

	return &orderv1.ListOrdersResponse{
		Orders:    converter.OrdersToProtoOrders(orders),
		TotalSize: int64(len(orders)),
	}, nil
}

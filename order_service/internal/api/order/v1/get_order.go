package v1

import (
	"context"

	orderv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/converter"
)

func (i *Implementation) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	order, err := i.orderService.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	return &orderv1.GetOrderResponse{
		Order: converter.OrderToProtoOrder(order),
	}, nil
}

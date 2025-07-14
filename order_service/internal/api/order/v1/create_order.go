package v1

import (
	"context"

	orderv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/converter"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	orderID, err := i.orderService.CreateOrder(ctx, converter.ProtoCreateOrderRequestToCreateOrderParams(req))
	if err != nil {
		return nil, err
	}

	return &orderv1.CreateOrderResponse{OrderId: orderID}, nil
}

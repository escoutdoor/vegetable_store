package v1

import (
	orderv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service"
)

type Implementation struct {
	orderService service.OrderService
	orderv1.UnimplementedOrderServiceServer
}

func NewImplementation(orderService service.OrderService) *Implementation {
	return &Implementation{
		orderService: orderService,
	}
}

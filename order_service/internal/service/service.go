package service

import (
	"context"

	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service/dto"
)

type OrderService interface {
	CreateOrder(ctx context.Context, in dto.CreateOrderParams) (string, error)
	GetOrder(ctx context.Context, orderID string) (entity.Order, error)
	ListOrders(ctx context.Context, in dto.ListOrdersParams) ([]entity.Order, error)
}

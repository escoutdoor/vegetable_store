package order

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service/dto"
)

func (s *service) ListOrders(ctx context.Context, in dto.ListOrdersParams) ([]entity.Order, error) {
	orders, err := s.orderRepository.ListOrders(ctx, in)
	if err != nil {
		return nil, errwrap.Wrap("get the list of orders from repository", err)
	}

	return orders, nil
}

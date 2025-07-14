package order

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
)

func (s *service) GetOrder(ctx context.Context, orderID string) (entity.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return entity.Order{}, errwrap.Wrap("get order from repository", err)
	}

	return order, nil
}

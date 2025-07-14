package order

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	vegetable_client "github.com/escoutdoor/vegetable_store/order_service/internal/client/vegetable"
	apperrors "github.com/escoutdoor/vegetable_store/order_service/internal/errors"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service/dto"
)

func (s *service) CreateOrder(ctx context.Context, in dto.CreateOrderParams) (string, error) {
	vegetableIDs := make([]string, 0, len(in.OrderItems))
	for _, oi := range in.OrderItems {
		vegetableIDs = append(vegetableIDs, oi.VegetableID)
	}

	list, err := s.vegetableClient.ListVegetables(ctx, newListVegetablesParams(vegetableIDs))
	if err != nil {
		return "", errwrap.Wrap("get the list of vegetables", err)
	}

	if len(list) != len(vegetableIDs) {
		var missingIDs []string
		for _, vid := range vegetableIDs {
			if _, ok := list[vid]; !ok {
				missingIDs = append(missingIDs, vid)
			}
		}
		if len(missingIDs) > 0 {
			return "", apperrors.VegetablesNotFound(missingIDs)
		}
	}

	orderItems := make([]dto.CreateOrderItemParams, 0, len(in.OrderItems))
	for _, oi := range in.OrderItems {
		orderItem := oi
		orderItem.DiscountedPrice = list[oi.VegetableID].DiscountedPrice
		orderItem.Price = list[oi.VegetableID].Price

		orderItems = append(orderItems, orderItem)
	}
	in.OrderItems = orderItems

	var orderID string
	if txErr := s.transactionManager.ReadCommited(ctx, func(ctx context.Context) error {
		in.TotalAmount = getTotalAmount(in)

		var err error
		orderID, err = s.orderRepository.CreateOrder(ctx, in)
		if err != nil {
			return errwrap.Wrap("create order", err)
		}

		updateVegetableParams := make([]vegetable_client.UpdateVegetableWeightParams, 0, len(in.OrderItems))
		for _, oi := range in.OrderItems {
			vegetable := list[oi.VegetableID]
			if vegetable.Weight < oi.Weight {
				return apperrors.InsufficientInventory(oi.VegetableID, vegetable.Weight, oi.Weight)
			}
			updateVegetableParams = append(updateVegetableParams, newUpdateVegetableWeightParams(vegetable.ID, vegetable.Weight-oi.Weight))
		}

		if err := s.vegetableClient.BatchUpdateVegetablesWeight(ctx, updateVegetableParams); err != nil {
			return errwrap.Wrap("batch update vegetables weight", err)
		}

		return nil
	}); txErr != nil {
		return "", txErr
	}

	return orderID, nil
}

func getTotalAmount(in dto.CreateOrderParams) float32 {
	var sum float32
	for _, oi := range in.OrderItems {
		if oi.DiscountedPrice != oi.Price {
			sum += oi.Weight * oi.DiscountedPrice
			continue
		}
		sum += oi.Weight * oi.Price
	}

	return sum
}

func newListVegetablesParams(vegetableIDs []string) vegetable_client.ListVegetablesParams {
	return vegetable_client.ListVegetablesParams{
		Limit:        len(vegetableIDs),
		VegetableIDs: vegetableIDs,
	}
}

func newUpdateVegetableWeightParams(vegetableID string, weight float32) vegetable_client.UpdateVegetableWeightParams {
	return vegetable_client.UpdateVegetableWeightParams{
		VegetableID: vegetableID,
		Weight:      weight,
	}
}

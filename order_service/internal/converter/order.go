package converter

import (
	orderv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/order/v1"
	"github.com/escoutdoor/vegetable_store/order_service/internal/entity"
	"github.com/escoutdoor/vegetable_store/order_service/internal/service/dto"
)

func ProtoListOrdersRequestToListOrdersParams(req *orderv1.ListOrdersRequest) dto.ListOrdersParams {
	return dto.ListOrdersParams{
		UserID:   req.UserId,
		Limit:    req.Limit,
		Offset:   req.Offset,
		OrderIDs: req.OrderIds,
	}
}

func ProtoCreateOrderRequestToCreateOrderParams(req *orderv1.CreateOrderRequest) dto.CreateOrderParams {
	orderItems := make([]dto.CreateOrderItemParams, 0, len(req.OrderItems))
	for _, oi := range req.OrderItems {
		orderItems = append(orderItems, protoOrderItemInputToCreateOrderItemParams(oi))
	}

	return dto.CreateOrderParams{
		UserID:     req.UserId,
		OrderItems: orderItems,
	}
}

func protoOrderItemInputToCreateOrderItemParams(orderItem *orderv1.OrderItemInput) dto.CreateOrderItemParams {
	return dto.CreateOrderItemParams{
		VegetableID: orderItem.VegetableId,
		Weight:      orderItem.Weight,

		Address: orderItem.AddressInfo.Address,

		FirstName:   orderItem.RecipientInfo.FirstName,
		LastName:    orderItem.RecipientInfo.LastName,
		PhoneNumber: orderItem.RecipientInfo.PhoneNumber,
		Email:       orderItem.RecipientInfo.Email,
	}
}

// func protoAddressToAddress(addressInfo *orderv1.AddressInfo) entity.AddressInfo {
// 	return entity.AddressInfo{
// 		Address: addressInfo.Address,
// 	}
// }
//
// func protoRecipientInfoToRecipientInfo(recipientInfo *orderv1.RecipientInfo) entity.RecipientInfo {
// 	return entity.RecipientInfo{
// 		FirstName:   recipientInfo.FirstName,
// 		LastName:    recipientInfo.LastName,
// 		PhoneNumber: recipientInfo.PhoneNumber,
// 		Email:       recipientInfo.Email,
// 	}
// }

func OrdersToProtoOrders(orders []entity.Order) []*orderv1.Order {
	list := make([]*orderv1.Order, 0, len(orders))
	for _, o := range orders {
		list = append(list, OrderToProtoOrder(o))
	}
	return list
}

func OrderToProtoOrder(order entity.Order) *orderv1.Order {
	orderItems := make([]*orderv1.OrderItem, 0, len(order.OrderItems))
	for _, io := range order.OrderItems {
		orderItems = append(orderItems, orderItemToProtoOrderItem(io))
	}

	return &orderv1.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		TotalAmount: order.TotalAmount,
		OrderItems:  orderItems,
	}
}

func orderItemToProtoOrderItem(orderItem entity.OrderItem) *orderv1.OrderItem {
	return &orderv1.OrderItem{
		OrderItemId:     orderItem.ID,
		VegetableId:     orderItem.VegetableID,
		Weight:          orderItem.Weight,
		Price:           orderItem.Price,
		DiscountedPrice: orderItem.DiscountedPrice,
		RecipientInfo:   recipientInfoToProtoRecipientInfo(orderItem.RecipientInfo),
		AddressInfo:     addressInfoToProtoAddressInfo(orderItem.AddressInfo),
	}
}

func recipientInfoToProtoRecipientInfo(recipientInfo entity.RecipientInfo) *orderv1.RecipientInfo {
	return &orderv1.RecipientInfo{
		FirstName:   recipientInfo.FirstName,
		LastName:    recipientInfo.LastName,
		PhoneNumber: recipientInfo.PhoneNumber,
		Email:       recipientInfo.Email,
	}
}

func addressInfoToProtoAddressInfo(addressInfo entity.AddressInfo) *orderv1.AddressInfo {
	return &orderv1.AddressInfo{
		Address: addressInfo.Address,
	}
}

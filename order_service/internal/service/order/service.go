package order

import (
	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/escoutdoor/vegetable_store/order_service/internal/client"
	"github.com/escoutdoor/vegetable_store/order_service/internal/repository"
)

type service struct {
	orderRepository    repository.OrderRepository
	transactionManager database.TxManager
	vegetableClient    client.VegetableClient
}

func NewService(
	orderRepository repository.OrderRepository,
	transactionManager database.TxManager,
	vegetableClient client.VegetableClient,
) *service {
	return &service{
		orderRepository:    orderRepository,
		transactionManager: transactionManager,
		vegetableClient:    vegetableClient,
	}
}

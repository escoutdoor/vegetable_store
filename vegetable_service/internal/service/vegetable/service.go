package vegetable

import (
	"github.com/escoutdoor/vegetable_store/common/pkg/database"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/repository"
)

type service struct {
	vegetableRepository repository.VegetableRepository
	transactionManager  database.TxManager
}

func NewService(vegetableRepository repository.VegetableRepository, transactionManager database.TxManager) *service {
	return &service{
		vegetableRepository: vegetableRepository,
		transactionManager:  transactionManager,
	}
}

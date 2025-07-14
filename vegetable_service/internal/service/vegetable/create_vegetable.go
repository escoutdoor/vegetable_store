package vegetable

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
)

func (s *service) CreateVegetable(ctx context.Context, in dto.CreateVegetableParams) (string, error) {
	id, err := s.vegetableRepository.CreateVegetable(ctx, in)
	if err != nil {
		return "", errwrap.Wrap("create vegetable", err)
	}

	return id, nil
}

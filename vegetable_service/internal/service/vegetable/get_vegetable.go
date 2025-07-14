package vegetable

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/entity"
)

func (s *service) GetVegetable(ctx context.Context, vegetableID string) (entity.Vegetable, error) {
	vegetable, err := s.vegetableRepository.GetVegetable(ctx, vegetableID)
	if err != nil {
		return entity.Vegetable{}, errwrap.Wrap("get vegetable from repository", err)
	}

	return vegetable, nil
}

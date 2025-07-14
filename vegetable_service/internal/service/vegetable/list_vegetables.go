package vegetable

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/entity"
)

func (s *service) ListVegetables(ctx context.Context, in dto.ListVegetablesParams) ([]entity.Vegetable, error) {
	list, err := s.vegetableRepository.ListVegetables(ctx, in)
	if err != nil {
		return nil, errwrap.Wrap("get a list of vegetables", err)
	}

	return list, nil
}

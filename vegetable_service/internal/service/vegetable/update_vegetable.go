package vegetable

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
)

func (s *service) UpdateVegetable(ctx context.Context, in dto.VegetableUpdateOperation) error {
	_, err := s.vegetableRepository.GetVegetable(ctx, in.ID)
	if err != nil {
		return errwrap.Wrap("get vegetable from repository", err)
	}

	if err := s.vegetableRepository.UpdateVegetable(ctx, in); err != nil {
		return errwrap.Wrap("update vegetable", err)
	}

	return nil
}

package vegetable

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
)

func (s *service) DeleteVegetable(ctx context.Context, vegetableID string) error {
	_, err := s.vegetableRepository.GetVegetable(ctx, vegetableID)
	if err != nil {
		return errwrap.Wrap("get vegetable from repository", err)
	}

	if err := s.vegetableRepository.DeleteVegetable(ctx, vegetableID); err != nil {
		return errwrap.Wrap("delete vegetable", err)
	}

	return nil
}

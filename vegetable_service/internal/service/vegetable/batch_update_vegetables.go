package vegetable

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/vegetable_service/internal/dto"
	apperrors "github.com/escoutdoor/vegetable_store/vegetable_service/internal/errors"
)

func (s *service) BatchUpdateVegetables(ctx context.Context, in dto.VegetablesBatchUpdateOperation) error {
	vegetableIDs := make([]string, 0, len(in.Requests))
	for _, r := range in.Requests {
		vegetableIDs = append(vegetableIDs, r.ID)
	}

	vegetables, err := s.vegetableRepository.ListVegetables(ctx, dto.ListVegetablesParams{
		Limit:        len(vegetableIDs),
		VegetableIDs: vegetableIDs,
	})
	if err != nil {
		return errwrap.Wrap("get the list of vegetables", err)
	}

	if len(vegetables) != len(vegetableIDs) {
		m := make(map[string]struct{}, len(vegetableIDs))
		for _, v := range vegetables {
			m[v.ID] = struct{}{}
		}

		missingIDs := make([]string, 0, (len(vegetableIDs) - len(vegetables)))
		for _, id := range vegetableIDs {
			if _, ok := m[id]; !ok {
				missingIDs = append(missingIDs, id)
			}
		}

		return apperrors.VegetablesNotFound(vegetableIDs)
	}

	if txErr := s.transactionManager.ReadCommited(ctx, func(ctx context.Context) error {
		for _, r := range in.Requests {
			if err := s.vegetableRepository.UpdateVegetable(ctx, r); err != nil {
				return errwrap.Wrap("update vegetable", err)
			}
		}

		return nil
	}); txErr != nil {
		return txErr
	}

	return nil
}

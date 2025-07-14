package user

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
)

func (s *service) ListUsers(ctx context.Context, in dto.ListUsersParams) ([]entity.User, error) {
	list, err := s.userRepository.ListUsers(ctx, in)
	if err != nil {
		return nil, errwrap.Wrap("list users from repository", err)
	}

	return list, nil
}

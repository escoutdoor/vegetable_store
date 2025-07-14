package user

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
)

func (s *service) UpdateUser(ctx context.Context, in dto.UserUpdateOperation) error {
	if _, err := s.userRepository.GetUserByID(ctx, in.ID); err != nil {
		return errwrap.Wrap("get user from repository", err)
	}

	if err := s.userRepository.UpdateUser(ctx, in); err != nil {
		return errwrap.Wrap("update user", err)
	}

	return nil
}

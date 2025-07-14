package user

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
)

func (s *service) DeleteUser(ctx context.Context, userID string) error {
	if _, err := s.userRepository.GetUserByID(ctx, userID); err != nil {
		return errwrap.Wrap("get user from repository", err)
	}

	if err := s.userRepository.DeleteUser(ctx, userID); err != nil {
		return errwrap.Wrap("delete user", err)
	}

	return nil
}

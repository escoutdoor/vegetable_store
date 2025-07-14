package user

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
)

func (s *service) GetUser(ctx context.Context, userID string) (entity.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return entity.User{}, errwrap.Wrap("get user from repository", err)
	}

	return user, nil
}

package auth

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
)

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (entity.Tokens, error) {
	userID, err := s.tokenProvider.ValidateRefreshToken(refreshToken)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("validate refresh token", err)
	}

	if _, err := s.userRepository.GetUserByID(ctx, userID); err != nil {
		return entity.Tokens{}, errwrap.Wrap("get user by if from repository", err)
	}

	tokens, err := s.tokenProvider.GenerateTokens(userID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate jwt tokens", err)
	}

	return tokens, nil
}

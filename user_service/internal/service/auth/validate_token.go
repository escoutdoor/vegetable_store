package auth

import (
	"context"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
)

func (s *service) ValidateToken(ctx context.Context, accessToken string) (string, error) {
	userID, err := s.tokenProvider.ValidateAccessToken(accessToken)
	if err != nil {
		return "", errwrap.Wrap("validate access token", err)
	}

	return userID, nil
}

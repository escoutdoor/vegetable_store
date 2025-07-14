package auth

import (
	"context"
	"errors"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
	apperrors "github.com/escoutdoor/vegetable_store/user_service/internal/errors"
	"github.com/escoutdoor/vegetable_store/user_service/internal/errors/codes"
	"github.com/escoutdoor/vegetable_store/user_service/internal/utils/hasher"
)

func (s *service) Login(ctx context.Context, email, password string) (entity.Tokens, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		appErr := new(apperrors.Error)
		if errors.As(err, &appErr) && appErr.Code == codes.UserNotFound {
			return entity.Tokens{}, apperrors.ErrIncorrectCreadentials
		}

		return entity.Tokens{}, errwrap.Wrap("get user by email from repository", err)
	}

	if match := hasher.CompareHashAndPassword(user.Password, password); !match {
		return entity.Tokens{}, apperrors.ErrIncorrectCreadentials
	}

	tokens, err := s.tokenProvider.GenerateTokens(user.ID)
	if err != nil {
		return entity.Tokens{}, errwrap.Wrap("generate jwt tokens", err)
	}

	return tokens, nil
}

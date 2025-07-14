package auth

import (
	"context"
	"errors"

	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
	apperrors "github.com/escoutdoor/vegetable_store/user_service/internal/errors"
	"github.com/escoutdoor/vegetable_store/user_service/internal/errors/codes"
	"github.com/escoutdoor/vegetable_store/user_service/internal/utils/hasher"
)

func (s *service) Register(ctx context.Context, in dto.CreateUserParams) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, in.Email)
	if err != nil {
		appErr := new(apperrors.Error)
		if errors.As(err, &appErr) {
			if appErr.Code != codes.UserNotFound {
				return "", errwrap.Wrap("get user by email", err)
			}
		} else {
			return "", errwrap.Wrap("get user by email", err)
		}
	}
	if user.Email != "" {
		return "", apperrors.EmailAlreadyExists(in.Email)
	}

	pw, err := hasher.HashPassword(in.Password)
	if err != nil {
		return "", errwrap.Wrap("hash password", err)
	}
	in.Password = pw

	userID, err := s.userRepository.CreateUser(ctx, in)
	if err != nil {
		return "", errwrap.Wrap("create user", err)
	}

	return userID, nil
}

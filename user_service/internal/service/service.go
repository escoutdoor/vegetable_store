package service

import (
	"context"

	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (entity.User, error)
	ListUsers(ctx context.Context, in dto.ListUsersParams) ([]entity.User, error)
	UpdateUser(ctx context.Context, in dto.UserUpdateOperation) error
	DeleteUser(ctx context.Context, userID string) error
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (entity.Tokens, error)
	Register(ctx context.Context, in dto.CreateUserParams) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (entity.Tokens, error)
	ValidateToken(ctx context.Context, accessToken string) (string, error)
}

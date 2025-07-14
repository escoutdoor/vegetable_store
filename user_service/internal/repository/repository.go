package repository

import (
	"context"

	"github.com/escoutdoor/vegetable_store/user_service/internal/dto"
	"github.com/escoutdoor/vegetable_store/user_service/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, in dto.CreateUserParams) (string, error)
	UpdateUser(ctx context.Context, in dto.UserUpdateOperation) error
	GetUserByID(ctx context.Context, userID string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	DeleteUser(ctx context.Context, userID string) error
	ListUsers(ctx context.Context, in dto.ListUsersParams) ([]entity.User, error)
}

package v1

import (
	"context"

	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/converter"
)

func (i *Implementation) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	panic("something happened to get user handler")
	user, err := i.userService.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{User: converter.UserToProtoUser(user)}, nil
}

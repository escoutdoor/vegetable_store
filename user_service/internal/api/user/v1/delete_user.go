package v1

import (
	"context"

	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
)

func (i *Implementation) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	err := i.userService.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &userv1.DeleteUserResponse{}, nil
}

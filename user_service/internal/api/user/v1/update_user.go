package v1

import (
	"context"

	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/converter"
)

func (i *Implementation) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	if err := validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	in, err := converter.ProtoUpdateUserRequestToUserUpdateOperation(req)
	if err != nil {
		return nil, err
	}

	if err := i.userService.UpdateUser(ctx, in); err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{}, nil
}

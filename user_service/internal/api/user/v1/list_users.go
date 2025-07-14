package v1

import (
	"context"

	userv1 "github.com/escoutdoor/vegetable_store/common/pkg/api/user/v1"
	"github.com/escoutdoor/vegetable_store/user_service/internal/converter"
)

func (i *Implementation) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	list, err := i.userService.ListUsers(ctx, converter.ProtoListUsersRequestToListUsersParams(req))
	if err != nil {
		return nil, err
	}

	return &userv1.ListUsersResponse{
		Users:     converter.UsersToProtoUsers(list),
		TotalSize: int64(len(list)),
	}, nil
}
